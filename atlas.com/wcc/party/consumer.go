package party

import (
	"atlas-wcc/kafka"
	"atlas-wcc/model"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameStatus       = "party_status_event"
	consumerNameMemberStatus = "party_member_status_event"
	topicTokenStatus         = "TOPIC_PARTY_STATUS"
	topicTokenMemberStatus   = "TOPIC_PARTY_MEMBER_STATUS"
)

func StatusConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[statusEvent](consumerNameStatus, topicTokenStatus, groupId, handleStatus(wid, cid))
	}
}

type statusEvent struct {
	WorldId     byte   `json:"world_id"`
	PartyId     uint32 `json:"party_id"`
	CharacterId uint32 `json:"character_id"`
	Type        string `json:"type"`
}

func handleStatus(wid byte, _ byte) kafka.HandlerFunc[statusEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event statusEvent) {
		if wid != event.WorldId {
			return
		}

		if event.Type == "CREATED" {
			session.ForSessionByCharacterId(event.CharacterId, showCreated(l, span)(event))
		} else if event.Type == "DISBANDED" {
			l.Debugf("Party %d disbanded.", event.PartyId)
		}
	}
}

func showCreated(l logrus.FieldLogger, _ opentracing.Span) func(event statusEvent) model.Operator[session.Model] {
	return func(event statusEvent) model.Operator[session.Model] {
		return func(s session.Model) error {
			l.Debugf("Party %d created for character %d.", event.PartyId, event.CharacterId)
			err := session.Announce(WritePartyCreated(l)(event.PartyId))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return err
		}
	}
}

func MemberStatusConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[memberStatusEvent](consumerNameMemberStatus, topicTokenMemberStatus, groupId, handleMemberStatus(wid, cid))
	}
}

type memberStatusEvent struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	PartyId     uint32 `json:"party_id"`
	CharacterId uint32 `json:"character_id"`
	Type        string `json:"type"`
}

func handleMemberStatus(wid byte, cid byte) kafka.HandlerFunc[memberStatusEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event memberStatusEvent) {
		if event.WorldId != wid && event.ChannelId != cid {
			return
		}
		if event.Type == "DISBANDED" {
			session.ForSessionByCharacterId(event.CharacterId, showDisbanded(l, span)(event.PartyId, event.CharacterId))
		}
	}
}

func showDisbanded(l logrus.FieldLogger, _ opentracing.Span) func(partyId uint32, characterId uint32) model.Operator[session.Model] {
	return func(partyId uint32, characterId uint32) model.Operator[session.Model] {
		return func(s session.Model) error {
			err := session.Announce(WritePartyDisbanded(l)(partyId, characterId))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return err
		}
	}
}
