package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type enableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func EnableActionsEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &enableActionsEvent{}
	}
}

func HandleEnableActionsEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*enableActionsEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, enableActions(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func enableActions(_ logrus.FieldLogger, _ *enableActionsEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteEnableActions())
	}
}
