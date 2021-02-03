package consumers

import (
	"atlas-wcc/domain"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func MapCharacterEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &mapCharacterEvent{}
	}
}

func HandleMapCharacterEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*mapCharacterEvent); ok {
			var handler processors.SessionOperator
			if event.Type == "ENTER" {
				handler = enterMap(l, *event)
			} else if event.Type == "EXIT" {
				handler = exitMap(l, *event)
			} else {
				l.Printf("[WARN] received a unhandled map character event type of %s.", event.Type)
				return
			}
			processors.ForSessionByCharacterId(event.CharacterId, handler)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleMapCharacterEvent]")
		}
	}
}

func enterMap(l *log.Logger, event mapCharacterEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		cIds, err := processors.GetCharacterIdsInMap(event.WorldId, event.ChannelId, event.MapId)
		if err != nil {
			return
		}

		cm := make(map[uint32]*domain.Character)
		for _, cId := range cIds {
			c, err := processors.GetCharacterById(cId)
			if err != nil {
				//log something
			} else {
				cm[c.Attributes().Id()] = c
			}
		}

		// Spawn new character for other character.
		for k, v := range cm {
			if k != event.CharacterId {
				s := *processors.GetSessionByCharacterId(k)
				s.Announce(writer.WriteSpawnCharacter(*v, *cm[event.CharacterId], true))
			}
		}

		// Spawn other characters for incoming character.
		for k, v := range cm {
			if k != event.CharacterId {
				session.Announce(writer.WriteSpawnCharacter(*cm[event.CharacterId], *v, false))
			}
		}

		// Spawn NPCs for incoming character.
		processors.ForEachNPCInMap(event.MapId, spawnNPCForSession(session))

		// Spawn monsters for incoming character.
		processors.ForEachMonsterInMap(event.WorldId, event.ChannelId, event.MapId, spawnMonsterForSession(session))

		// Spawn drops for incoming character.
		processors.ForEachDropInMap(event.WorldId, event.ChannelId, event.MapId, spawnDropForSession(session))
	}
}

func spawnDropForSession(session mapleSession.MapleSession) processors.DropOperator {
	return func(drop domain.Drop) {
		var a = uint32(0)
		if drop.ItemId() != 0 {
			a = 0
		} else {
			a = drop.Meso()
		}
		session.Announce(writer.WriteDropItemFromMapObject(drop.UniqueId(), drop.ItemId(), drop.Meso(), a,
			drop.DropperUniqueId(), drop.DropType(), drop.OwnerId(), drop.OwnerPartyId(), session.CharacterId(),
			0, drop.DropTime(), drop.DropX(), drop.DropY(), drop.DropperX(), drop.DropperY(),
			drop.CharacterDrop(), drop.Mod()))
	}
}

func spawnMonsterForSession(session mapleSession.MapleSession) processors.MonsterOperator {
	return func(monster domain.Monster) {
		session.Announce(writer.WriteSpawnMonster(monster, false))
	}
}

func spawnNPCForSession(session mapleSession.MapleSession) processors.NPCOperator {
	return func(npc domain.NPC) {
		session.Announce(writer.WriteSpawnNPC(npc))
		session.Announce(writer.WriteSpawnNPCController(npc, true))
	}
}

func exitMap(_ *log.Logger, event mapCharacterEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		processors.ForEachOtherSessionInMap(event.WorldId, event.ChannelId, event.CharacterId, removeCharacterForSession(event.CharacterId))
		processors.ForEachNPCInMap(event.MapId, removeNpcForSession(session))
	}
}

func removeNpcForSession(session mapleSession.MapleSession) processors.NPCOperator {
	return func(npc domain.NPC) {
		session.Announce(writer.WriteRemoveNPCController(npc.ObjectId()))
		session.Announce(writer.WriteRemoveNPC(npc.ObjectId()))
	}
}

func removeCharacterForSession(characterId uint32) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteRemoveCharacterFromMap(characterId))
	}
}
