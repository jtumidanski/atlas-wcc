package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/monster"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type monsterEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
	MonsterId uint32 `json:"monsterId"`
	ActorId   uint32 `json:"actorId"`
	Type      string `json:"type"`
}

func MonsterEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &monsterEvent{}
	}
}

func HandleMonsterEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*monsterEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			m, err := monster.GetById(l, span)(event.UniqueId)
			if err != nil {
				l.WithError(err).Errorf("Unable to monster %d to create.", event.UniqueId)
				return
			}

			var h session.Operator
			if event.Type == "CREATED" {
				h = createMonster(l, event, m)
			} else if event.Type == "DESTROYED" {
				h = destroyMonster(l, event)
			} else {
				l.Warnf("Unable to handle %s event type for monster events.", event.Type)
				return
			}

			session.ForEachInMap(l, span)(wid, cid, event.MapId, h)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func destroyMonster(l logrus.FieldLogger, event *monsterEvent) session.Operator {
	k1 := writer.WriteKillMonster(l)(event.UniqueId, false)
	k2 := writer.WriteKillMonster(l)(event.UniqueId, true)
	return func(s *session.Model) {
		err := s.Announce(k1)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		err = s.Announce(k2)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func createMonster(l logrus.FieldLogger, _ *monsterEvent, monster *monster.Model) session.Operator {
	sm := writer.WriteSpawnMonster(l)(monster, false)
	return func(s *session.Model) {
		err := s.Announce(sm)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
