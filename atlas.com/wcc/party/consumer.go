package party

import (
	"atlas-wcc/kafka"
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
			l.Debugf("Party %d created for character %d.", event.PartyId, event.CharacterId)
			session.IfPresentByCharacterId(event.CharacterId, session.AnnounceOperator(WritePartyCreated(l)(event.PartyId)))
		} else if event.Type == "DISBANDED" {
			l.Debugf("Party %d disbanded.", event.PartyId)
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
			session.IfPresentByCharacterId(event.CharacterId, session.AnnounceOperator(WritePartyDisbanded(l)(event.PartyId, event.CharacterId)))
		}
	}
}
