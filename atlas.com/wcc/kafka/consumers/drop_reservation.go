package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type dropReservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

func DropReservationEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &dropReservationEvent{}
	}
}

func HandleDropReservationEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*dropReservationEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			if event.Type == "SUCCESS" {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, cancelDropReservation(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func cancelDropReservation(_ logrus.FieldLogger, _ *dropReservationEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteEnableActions())
	}
}
