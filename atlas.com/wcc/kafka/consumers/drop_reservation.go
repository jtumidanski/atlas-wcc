package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type dropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

func DropReservationEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &dropReservationEvent{}
	}
}

func HandleDropReservationEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*dropReservationEvent); ok {
			if event.Type == "SUCCESS" {
				return
			}

			processors.ForSessionByCharacterId(l, event.CharacterId, cancelDropReservation(event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleDropReservationEvent]")
		}
	}
}

func cancelDropReservation(_ *dropReservationEvent) processors.SessionOperator {
	return func(l *log.Logger, session mapleSession.MapleSession) {
		session.Announce(writer.WriteEnableActions())
	}
}
