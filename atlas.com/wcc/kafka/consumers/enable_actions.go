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
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, enableActions(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func enableActions(_ logrus.FieldLogger, _ *enableActionsEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteEnableActions())
	}
}
