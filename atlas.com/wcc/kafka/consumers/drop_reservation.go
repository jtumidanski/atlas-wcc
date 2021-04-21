package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*dropReservationEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			if event.Type == "SUCCESS" {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, cancelDropReservation(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func cancelDropReservation(_ logrus.FieldLogger, _ *dropReservationEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteEnableActions())
	}
}
