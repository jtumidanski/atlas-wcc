package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"fmt"
	"github.com/sirupsen/logrus"
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

type ServerNoticeEvent struct {
	Type        string `json:"type"`
	RecipientId uint32 `json:"recipientId"`
	Message     string `json:"message"`
}

func ServerNoticeEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &ServerNoticeEvent{}
	}
}

func HandleServerNoticeEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*ServerNoticeEvent); ok {
			if actingSession := session.GetByCharacterId(event.RecipientId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.RecipientId, showServerNotice(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showServerNotice(l logrus.FieldLogger, event *ServerNoticeEvent) session.Operator {
	return func(s *session.Model) {
		err := s.Announce(writer.WriteServerNotice(l)(s.ChannelId(), getServerNoticeByType(event.Type), event.Message, false, 0))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func getServerNoticeByType(t string) byte {
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
