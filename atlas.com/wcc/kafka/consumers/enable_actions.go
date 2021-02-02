package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
)

type enableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func EnableActionsEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &enableActionsEvent{}
	}
}

func HandleEnableActionsEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*enableActionsEvent); ok {
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}
			(*as).Announce(writer.WriteEnableActions())
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleEnableActionsEvent]")
		}
	}
}
