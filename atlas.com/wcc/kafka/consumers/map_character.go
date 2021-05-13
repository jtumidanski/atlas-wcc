package consumers

import (
	"atlas-wcc/domain"
	"atlas-wcc/kafka/handler"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type mapCharacterEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func MapCharacterEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &mapCharacterEvent{}
	}
}

func HandleMapCharacterEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*mapCharacterEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			if event.Type == "ENTER" {
				processors.ForSessionByCharacterId(event.CharacterId, enterMap(l, *event))
			} else if event.Type == "EXIT" {
				processors.ForEachOtherSessionInMap(event.WorldId, event.ChannelId, event.CharacterId, removeCharacterForSession(event.CharacterId))
			} else {
				l.Warnf("Received a unhandled map character event type of %s.", event.Type)
				return
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func enterMap(l logrus.FieldLogger, event mapCharacterEvent) processors.SessionOperator {
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

func removeCharacterForSession(characterId uint32) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteRemoveCharacterFromMap(characterId))
	}
}
