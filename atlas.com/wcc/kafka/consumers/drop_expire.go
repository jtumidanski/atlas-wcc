package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
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
			processors.ForEachSessionInMap(l, wid, cid, event.MapId, expireItem(event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleDropExpireEvent]")
		}
	}
}

func expireItem(event *DropExpireEvent) processors.SessionOperator {
	return func(l *log.Logger, session mapleSession.MapleSession) {
		session.Announce(writer.WriteRemoveItem(event.UniqueId, 0, 0))
	}
}
