package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type DropExpireEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
}

func DropExpireEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &DropExpireEvent{}
	}
}

func HandleDropExpireEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*DropExpireEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			session.ForEachSessionInMap(wid, cid, event.MapId, expireItem(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func expireItem(_ logrus.FieldLogger, event *DropExpireEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteRemoveItem(event.UniqueId, 0, 0))
	}
}
