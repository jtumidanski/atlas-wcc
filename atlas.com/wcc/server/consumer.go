package server

import (
	"atlas-wcc/kafka"
	"atlas-wcc/session"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameNotice = "server_notice_command"
	topicTokenNotice   = "TOPIC_SERVER_NOTICE_COMMAND"
)

const (
	notice         byte = 0
	popUp          byte = 1
	megaphone      byte = 2
	superMegaphone byte = 3
	scroll         byte = 4
	pinkText       byte = 5
	lightBlue      byte = 6
)

func NoticeConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[noticeEvent](consumerNameNotice, topicTokenNotice, groupId, handleNotice(wid, cid))
	}
}

type noticeEvent struct {
	Type        string `json:"type"`
	RecipientId uint32 `json:"recipientId"`
	Message     string `json:"message"`
}

func handleNotice(_ byte, _ byte) kafka.HandlerFunc[noticeEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event noticeEvent) {
		if actingSession := session.GetByCharacterId(event.RecipientId); actingSession == nil {
			return
		}

		session.ForSessionByCharacterId(event.RecipientId, showNotice(l, event))
	}
}

func showNotice(l logrus.FieldLogger, event noticeEvent) session.Operator {
	return func(s *session.Model) {
		err := s.Announce(WriteServerNotice(l)(s.ChannelId(), getNoticeByType(event.Type), event.Message, false, 0))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func getNoticeByType(t string) byte {
	switch t {
	case "NOTICE":
		return notice
	case "POP_UP":
		return popUp
	case "MEGAPHONE":
		return megaphone
	case "SUPER_MEGAPHONE":
		return superMegaphone
	case "SCROLL":
		return scroll
	case "PINK_TEXT":
		return pinkText
	case "LIGHT_BLUE":
		return lightBlue
	}
	panic(fmt.Sprintf("unsupported server notice type %s", t))
}
