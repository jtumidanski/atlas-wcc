package consumers

import (
	"atlas-wcc/character"
	"atlas-wcc/drop"
	"atlas-wcc/kafka/handler"
	_map "atlas-wcc/map"
	"atlas-wcc/monster"
	"atlas-wcc/npc"
	"atlas-wcc/reactor"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
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
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*mapCharacterEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			if event.Type == "ENTER" {
				session.ForSessionByCharacterId(event.CharacterId, enterMap(l, span)(*event))
			} else if event.Type == "EXIT" {
				session.ForEachOtherInMap(l, span)(event.WorldId, event.ChannelId, event.CharacterId, removeCharacterForSession(l)(event.CharacterId))
			} else {
				l.Warnf("Received a unhandled map character event type of %s.", event.Type)
				return
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func enterMap(l logrus.FieldLogger, span opentracing.Span) func(event mapCharacterEvent) session.Operator {
	return func(event mapCharacterEvent) session.Operator {

		return func(s *session.Model) {
			cIds, err := _map.GetCharacterIdsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId)
			if err != nil {
				return
			}

			cm := make(map[uint32]*character.Model)
			for _, cId := range cIds {
				c, err := character.GetCharacterById(l, span)(cId)
				if err != nil {
					//log something
				} else {
					cm[c.Attributes().Id()] = c
				}
			}

			// Spawn new character for other character.
			for k, v := range cm {
				if k != event.CharacterId {
					as := *session.GetByCharacterId(k)
					err = as.Announce(writer.WriteSpawnCharacter(l)(*v, *cm[event.CharacterId], true))
					if err != nil {
						l.WithError(err).Errorf("Unable to spawn character %d for %d", event.CharacterId, v.Attributes().Id())
					}
				}
			}

			// Spawn other characters for incoming character.
			for k, v := range cm {
				if k != event.CharacterId {
					err = s.Announce(writer.WriteSpawnCharacter(l)(*cm[event.CharacterId], *v, false))
					if err != nil {
						l.WithError(err).Errorf("Unable to spawn character %d for %d", v.Attributes().Id(), event.CharacterId)
					}
				}
			}

			// Spawn NPCs for incoming character.
			npc.ForEachInMap(l, span)(event.MapId, spawnNPCForSession(l)(s))

			// Spawn monsters for incoming character.
			monster.ForEachInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, spawnMonsterForSession(l)(s))

			// Spawn drops for incoming character.
			drop.ForEachInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, spawnDropForSession(l)(s))

			// Spawn reactors for incoming character.
			reactor.ForEachAliveInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, spawnReactorForSession(l)(s))
		}
	}
}

func spawnReactorForSession(l logrus.FieldLogger) func(s *session.Model) reactor.ModelOperator {
	return func(s *session.Model) reactor.ModelOperator {
		return func(r *reactor.Model) {
			err := s.Announce(writer.WriteReactorSpawn(l)(r.Id(), r.Classification(), r.State(), r.X(), r.Y()))
			if err != nil {
				l.WithError(err).Errorf("Unable to show reactor %d creation to session %d.", r.Id(), s.SessionId())
			}
		}
	}
}

func spawnDropForSession(l logrus.FieldLogger) func(s *session.Model) drop.Operator {
	return func(s *session.Model) drop.Operator {
		return func(drop *drop.Model) {
			var a = uint32(0)
			if drop.ItemId() != 0 {
				a = 0
			} else {
				a = drop.Meso()
			}
			err := s.Announce(writer.WriteDropItemFromMapObject(l)(drop.UniqueId(), drop.ItemId(), drop.Meso(), a,
				drop.DropperUniqueId(), drop.DropType(), drop.OwnerId(), drop.OwnerPartyId(), s.CharacterId(),
				0, drop.DropTime(), drop.DropX(), drop.DropY(), drop.DropperX(), drop.DropperY(),
				drop.CharacterDrop(), drop.Mod()))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce drop to character %d", s.CharacterId())
			}
		}
	}
}

func spawnMonsterForSession(l logrus.FieldLogger) func(s *session.Model) monster.ModelOperator {
	return func(s *session.Model) monster.ModelOperator {
		return func(monster *monster.Model) {
			err := s.Announce(writer.WriteSpawnMonster(l)(monster, false))
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn monster %d for character %d", monster.MonsterId(), s.CharacterId())
			}
		}
	}
}

func spawnNPCForSession(l logrus.FieldLogger) func(s *session.Model) npc.ModelOperator {
	return func(s *session.Model) npc.ModelOperator {
		return func(npc *npc.Model) {
			err := s.Announce(writer.WriteSpawnNPC(l)(npc))
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn npc %d for character %d", npc.Id(), s.CharacterId())
			}
			err = s.Announce(writer.WriteSpawnNPCController(l)(npc, true))
			if err != nil {
				l.WithError(err).Errorf("Unable to spawn npc controller %d for character %d", npc.Id(), s.CharacterId())
			}
		}
	}
}

func removeCharacterForSession(l logrus.FieldLogger) func(characterId uint32) session.Operator {
	return func(characterId uint32) session.Operator {
		return func(s *session.Model) {
			err := s.Announce(writer.WriteRemoveCharacterFromMap(l)(characterId))
			if err != nil {
				l.WithError(err).Errorf("Unable to remove character %d from view for character %d", characterId, s.CharacterId())
			}
		}
	}
}
