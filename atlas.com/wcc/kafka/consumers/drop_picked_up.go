package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
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
			processors.ForEachSessionInMap(l, wid, cid, event.MapId, removeItem(event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleDropPickedUpEvent]")
		}
	}
}

func removeItem(event *dropPickedUpEvent) processors.SessionOperator {
	return func(l *log.Logger, session mapleSession.MapleSession) {
		session.Announce(writer.WriteRemoveItem(event.DropId, 2, event.CharacterId))
	}
}
