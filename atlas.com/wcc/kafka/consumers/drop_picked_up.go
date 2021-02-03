package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
)

type dropPickedUpEvent struct {
	DropId      uint32 `json:"dropId"`
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
}

func DropPickedUpEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &dropPickedUpEvent{}
	}
}

func HandleDropPickedUpEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*dropPickedUpEvent); ok {
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}
			sl, err := getSessionsForThoseInMap(wid, cid, event.MapId)
			if err != nil {
				return
			}
			for _, s := range sl {
				s.Announce(writer.WriteRemoveItem(event.DropId, 2, event.CharacterId))
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleDropPickedUpEvent]")
		}
	}
}
