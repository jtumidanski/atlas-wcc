package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
)

type EnableActionsEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func EnableActionsEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &EnableActionsEvent{}
	}
}

func HandleEnableActionsEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, event interface{}) {
		e := *event.(*EnableActionsEvent)
		as := getSessionForCharacterId(e.CharacterId)
		if as == nil {
			l.Printf("[ERROR] unable to locate session for character %d.", e.CharacterId)
			return
		}

		(*as).Announce(writer.WriteEnableActions())
	}
}
