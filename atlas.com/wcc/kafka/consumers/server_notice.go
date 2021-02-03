package consumers

import (
	"atlas-wcc/socket/response/writer"
	"fmt"
	"log"
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

func ServerNoticeEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &ServerNoticeEvent{}
	}
}

func HandleServerNoticeEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*ServerNoticeEvent); ok {
			as := getSessionForCharacterId(event.RecipientId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.RecipientId)
				return
			}
			(*as).Announce(writer.WriteServerNotice((*as).ChannelId(), getServerNoticeByType(event.Type), event.Message, false, 0))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleServerNoticeEvent]")
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