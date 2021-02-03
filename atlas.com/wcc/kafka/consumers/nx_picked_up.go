package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"fmt"
	"log"
)

const nxGainFormat = "You have earned #e#b%d NX#k#n."

type nxPickedUpEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

func NXPickedUpEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &nxPickedUpEvent{}
	}
}

func HandleNXPickedUpEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*nxPickedUpEvent); ok {
			processors.ForSessionByCharacterId(event.CharacterId, showNXGain(l, event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleNXPickedUpEvent]")
		}
	}
}

func showNXGain(_ *log.Logger, event *nxPickedUpEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteHint(fmt.Sprintf(nxGainFormat, event.Gain), 300, 10))
		session.Announce(writer.WriteEnableActions())
	}
}
