package consumers

import (
	"atlas-wcc/character"
	"atlas-wcc/drop"
	"atlas-wcc/kafka/handler"
	"atlas-wcc/monster"
	"atlas-wcc/npc"
	"atlas-wcc/session"
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
				session.ForSessionByCharacterId(event.CharacterId, enterMap(l, *event))
			} else if event.Type == "EXIT" {
				session.ForEachOtherSessionInMap(event.WorldId, event.ChannelId, event.CharacterId, removeCharacterForSession(event.CharacterId))
			} else {
				l.Warnf("Received a unhandled map character event type of %s.", event.Type)
				return
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func enterMap(l logrus.FieldLogger, event mapCharacterEvent) session.SessionOperator {
	return func(s session.Model) {
		cIds, err := character.GetCharacterIdsInMap(event.WorldId, event.ChannelId, event.MapId)
		if err != nil {
			return
		}

		cm := make(map[uint32]*character.Model)
		for _, cId := range cIds {
			c, err := character.GetCharacterById(cId)
			if err != nil {
				//log something
			} else {
				cm[c.Attributes().Id()] = c
			}
		}

		// Spawn new character for other character.
		for k, v := range cm {
			if k != event.CharacterId {
				s := *session.GetSessionByCharacterId(k)
				s.Announce(writer.WriteSpawnCharacter(*v, *cm[event.CharacterId], true))
			}
		}

		// Spawn other characters for incoming character.
		for k, v := range cm {
			if k != event.CharacterId {
				s.Announce(writer.WriteSpawnCharacter(*cm[event.CharacterId], *v, false))
			}
		}

		// Spawn NPCs for incoming character.
		npc.ForEachNPCInMap(event.MapId, spawnNPCForSession(s))

		// Spawn monsters for incoming character.
		monster.ForEachMonsterInMap(event.WorldId, event.ChannelId, event.MapId, spawnMonsterForSession(s))

		// Spawn drops for incoming character.
		drop.ForEachDropInMap(event.WorldId, event.ChannelId, event.MapId, spawnDropForSession(s))
	}
}

func spawnDropForSession(s session.Model) drop.DropOperator {
	return func(drop drop.Model) {
		var a = uint32(0)
		if drop.ItemId() != 0 {
			a = 0
		} else {
			a = drop.Meso()
		}
		s.Announce(writer.WriteDropItemFromMapObject(drop.UniqueId(), drop.ItemId(), drop.Meso(), a,
			drop.DropperUniqueId(), drop.DropType(), drop.OwnerId(), drop.OwnerPartyId(), s.CharacterId(),
			0, drop.DropTime(), drop.DropX(), drop.DropY(), drop.DropperX(), drop.DropperY(),
			drop.CharacterDrop(), drop.Mod()))
	}
}

func spawnMonsterForSession(s session.Model) monster.MonsterOperator {
	return func(monster monster.Model) {
		s.Announce(writer.WriteSpawnMonster(monster, false))
	}
}

func spawnNPCForSession(s session.Model) npc.NPCOperator {
	return func(npc npc.Model) {
		s.Announce(writer.WriteSpawnNPC(npc))
		s.Announce(writer.WriteSpawnNPCController(npc, true))
	}
}

func removeCharacterForSession(characterId uint32) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteRemoveCharacterFromMap(characterId))
	}
}
