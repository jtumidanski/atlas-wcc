package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type partyMemberStatusEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	PartyId     uint32 `json:"party_id"`
	CharacterId uint32 `json:"character_id"`
	Type        string `json:"type"`
}

func EmptyPartyMemberStatusEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &partyMemberStatusEvent{}
	}
}

func HandlePartyMemberStatusEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*partyMemberStatusEvent); ok {
			if event.WorldId != wid && event.ChannelId != cid {
				return
			}
			if event.Type == "DISBANDED" {
				session.ForSessionByCharacterId(event.CharacterId, partyDisbanded(l, span)(event.PartyId, event.CharacterId))
			}
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func partyDisbanded(l logrus.FieldLogger, _ opentracing.Span) func(partyId uint32, characterId uint32) session.Operator {
	return func(partyId uint32, characterId uint32) session.Operator {
		return func(model *session.Model) {
			err := model.Announce(writer.WritePartyDisbanded(l)(partyId, characterId))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", model.CharacterId())
			}
		}
	}
}