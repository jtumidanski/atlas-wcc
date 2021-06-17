package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/monster"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type monsterControlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	UniqueId    uint32 `json:"uniqueId"`
	Type        string `json:"type"`
}

func MonsterControlEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &monsterControlEvent{}
	}
}

func HandleMonsterControlEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*monsterControlEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			var h session.Operator
			if event.Type == "START" {
				h = startControl(l, event)
			} else if event.Type == "STOP" {
				h = stopControl(l, event)
			} else {
				l.Warnf("Received unhandled monster control event type of %s", event.Type)
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, h)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func stopControl(l logrus.FieldLogger, event *monsterControlEvent) session.Operator {
	return func(s *session.Model) {
		m, err := monster.GetById(event.UniqueId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve monster %d for control change", event.UniqueId)
			return
		}
		l.Infof("Stopping control of %d for character %d.", event.UniqueId, event.CharacterId)
		err = s.Announce(writer.WriteStopControlMonster(l)(m))
		if err != nil {
			l.WithError(err).Errorf("Unable to stop control of %d by %d", event.UniqueId, event.CharacterId)
		}
	}
}

func startControl(l logrus.FieldLogger, event *monsterControlEvent) session.Operator {
	return func(s *session.Model) {
		m, err := monster.GetById(event.UniqueId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve monster %d for control change", event.UniqueId)
			return
		}
		l.Infof("Starting control of %d for character %d.", event.UniqueId, event.CharacterId)
		err = s.Announce(writer.WriteControlMonster(l)(m, false, false))
		if err != nil {
			l.WithError(err).Errorf("Unable to start control of %d by %d", event.UniqueId, event.CharacterId)
		}
	}
}
