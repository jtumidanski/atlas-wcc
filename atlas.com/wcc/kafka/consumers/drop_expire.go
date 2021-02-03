package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
)

type DropExpireEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
}

func DropExpireEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &DropExpireEvent{}
	}
}

func HandleDropExpireEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*DropExpireEvent); ok {
			sl, err := getSessionsForThoseInMap(wid, cid, event.MapId)
			if err != nil {
				return
			}
			for _, s := range sl {
				s.Announce(writer.WriteRemoveItem(event.UniqueId, 0, 0))
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleDropExpireEvent]")
		}
	}
}
