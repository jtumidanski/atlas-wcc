package wishlist

import (
	"atlas-wcc/kafka"
	"atlas-wcc/model"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameStatus = "wishlist_status_event"
	topicTokenStatus   = "TOPIC_WISHLIST_STATUS_EVENT"
)

func StatusEventConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[statusEvent](consumerNameStatus, topicTokenStatus, groupId, handleStatusEvent(wid, cid))
	}
}

type statusEvent struct {
	CharacterId uint32 `json:"character_id"`
	Status      string `json:"status"`
}

func handleStatusEvent(_ byte, _ byte) kafka.HandlerFunc[statusEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event statusEvent) {
		session.ForSessionByCharacterId(event.CharacterId, showWishList(l, span)(event.CharacterId))
	}
}

func showWishList(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) model.Operator[session.Model] {
	return func(characterId uint32) model.Operator[session.Model] {
		wl, err := GetById(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve wishlist for character %d.", characterId)
			return model.ErrorOperator[session.Model](err)
		}
		return session.Announce(WriteWishList(l)(wl, true))
	}
}
