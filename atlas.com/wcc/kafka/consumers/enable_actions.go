package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
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
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*enableActionsEvent); ok {
			processors.ForSessionByCharacterId(event.CharacterId, enableActions(l, event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleEnableActionsEvent]")
		}
	}
}

func enableActions(_ *log.Logger, _ *enableActionsEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteEnableActions())
	}
}
