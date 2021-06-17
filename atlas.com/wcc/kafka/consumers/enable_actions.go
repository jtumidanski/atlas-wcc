package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type enableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func EnableActionsEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &enableActionsEvent{}
	}
}

func HandleEnableActionsEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*enableActionsEvent); ok {
			if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, enableActions(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func enableActions(l logrus.FieldLogger, _ *enableActionsEvent) session.Operator {
	b := writer.WriteEnableActions(l)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
