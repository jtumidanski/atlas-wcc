package wishlist

import (
	"atlas-wcc/kafka"
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
		s, err := session.GetByCharacterId(event.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate session for character %d entering cash shop.", event.CharacterId)
			return
		}

		wl, err := GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve wishlist for character %d.", event.CharacterId)
		}
		err = session.Announce(WriteWishList(l)(wl, true))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write wish list to character %d.", event.CharacterId)
			return
		}
	}
}
