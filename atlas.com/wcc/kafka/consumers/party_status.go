package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type partyStatusEvent struct {
	WorldId     byte   `json:"world_id"`
	PartyId     uint32 `json:"party_id"`
	CharacterId uint32 `json:"character_id"`
	Type        string `json:"type"`
}

func EmptyPartyStatusEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &partyStatusEvent{}
	}
}

func HandlePartyStatusEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*partyStatusEvent); ok {
			if wid != event.WorldId {
				return
			}

			if event.Type == "CREATED" {
				session.ForSessionByCharacterId(event.CharacterId, partyCreated(l, span)(event))
			} else if event.Type == "DISBANDED" {
				l.Debugf("Party %d disbanded.", event.PartyId)
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func partyCreated(l logrus.FieldLogger, _ opentracing.Span) func(event *partyStatusEvent) session.Operator {
	return func(event *partyStatusEvent) session.Operator {
		return func(model *session.Model) {
			l.Debugf("Party %d created for character %d.", event.PartyId, event.CharacterId)
			err := model.Announce(writer.WritePartyCreated(l)(event.PartyId))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", model.CharacterId())
			}
		}
	}
}
