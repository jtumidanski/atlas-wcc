package monster

import (
	"atlas-wcc/kafka"
	"atlas-wcc/model"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameControl = "control_monster_event"
	topicTokenControl   = "TOPIC_CONTROL_MONSTER_EVENT"
)

func ControlConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[controlEvent](consumerNameControl, topicTokenControl, groupId, handleControl(wid, cid))
	}
}

type controlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	UniqueId    uint32 `json:"uniqueId"`
	Type        string `json:"type"`
}

func handleControl(wid byte, cid byte) kafka.HandlerFunc[controlEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event controlEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		var h model.Operator[session.Model]
		if event.Type == "START" {
			h = startControl(l, span, event)
		} else if event.Type == "STOP" {
			h = stopControl(l, span, event)
		} else {
			l.Warnf("Received unhandled monster control event type of %s", event.Type)
			return
		}

		session.IfPresentByCharacterId(event.CharacterId, h)
	}
}

func stopControl(l logrus.FieldLogger, span opentracing.Span, event controlEvent) model.Operator[session.Model] {
	m, err := GetById(l, span)(event.UniqueId)
	if err != nil {
		l.WithError(err).Errorf("Unable to retrieve monster %d for control change", event.UniqueId)
		return model.ErrorOperator[session.Model](err)
	}
	l.Infof("Stopping control of %d for character %d.", event.UniqueId, event.CharacterId)
	return session.AnnounceOperator(WriteStopControlMonster(l)(m))
}

func startControl(l logrus.FieldLogger, span opentracing.Span, event controlEvent) model.Operator[session.Model] {
	m, err := GetById(l, span)(event.UniqueId)
	if err != nil {
		l.WithError(err).Errorf("Unable to retrieve monster %d for control change", event.UniqueId)
		return model.ErrorOperator[session.Model](err)
	}
	return session.AnnounceOperator(WriteControlMonster(l)(m, false, false))
}
