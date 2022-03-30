package drop

import (
	"atlas-wcc/character"
	"atlas-wcc/character/properties"
	"atlas-wcc/kafka"
	"atlas-wcc/session"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameReservation = "drop_reservation_event"
	consumerNamePickupItem  = "picked_up_item_event"
	consumerNamePickupNX    = "picked_up_nx_event"
	topicTokenPickupItem    = "TOPIC_PICKED_UP_ITEM"
	topicTokenPickupNX      = "TOPIC_PICKED_UP_NX"
	topicTokenReservation   = "TOPIC_DROP_RESERVATION_EVENT"
)

func ReservationConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[reservationEvent](consumerNameReservation, topicTokenReservation, groupId, handleReservation(wid, cid))
	}
}

type reservationEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
	Type        string `json:"type"`
}

func handleReservation(_ byte, _ byte) kafka.HandlerFunc[reservationEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event reservationEvent) {
		if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
			return
		}

		if event.Type == "SUCCESS" {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, cancelDropReservation(l, event))
	}
}

func cancelDropReservation(l logrus.FieldLogger, _ reservationEvent) session.Operator {
	b := properties.WriteEnableActions(l)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func PickupItemConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[pickupItemEvent](consumerNamePickupItem, topicTokenPickupItem, groupId, handlePickupItem(wid, cid))
	}
}

type pickupItemEvent struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    uint32 `json:"quantity"`
}

func handlePickupItem(_ byte, _ byte) kafka.HandlerFunc[pickupItemEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event pickupItemEvent) {
		if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, showItemGain(l, event))
	}
}

func showItemGain(l logrus.FieldLogger, event pickupItemEvent) session.Operator {
	ig := properties.WriteShowItemGain(l)(event.ItemId, event.Quantity)
	ea := properties.WriteEnableActions(l)
	return func(s *session.Model) {
		err := s.Announce(ig)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		err = s.Announce(ea)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func PickupNXConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[pickupNXEvent](consumerNamePickupNX, topicTokenPickupNX, groupId, handlePickupNX(wid, cid))
	}
}

const nxGainFormat = "You have earned #e#b%d NX#k#n."

type pickupNXEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

func handlePickupNX(_ byte, _ byte) kafka.HandlerFunc[pickupNXEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event pickupNXEvent) {
		if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, showNXGain(l, event))
	}
}

func showNXGain(l logrus.FieldLogger, event pickupNXEvent) session.Operator {
	h := character.WriteHint(l)(fmt.Sprintf(nxGainFormat, event.Gain), 300, 10)
	ea := properties.WriteEnableActions(l)
	return func(s *session.Model) {
		err := s.Announce(h)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		err = s.Announce(ea)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
