package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type DropEvent struct {
	WorldId         byte   `json:"worldId"`
	ChannelId       byte   `json:"channelId"`
	MapId           uint32 `json:"mapId"`
	UniqueId        uint32 `json:"uniqueId"`
	ItemId          uint32 `json:"itemId"`
	Quantity        uint32 `json:"quantity"`
	Meso            uint32 `json:"meso"`
	DropType        byte   `json:"dropType"`
	DropX           int16  `json:"dropX"`
	DropY           int16  `json:"dropY"`
	OwnerId         uint32 `json:"ownerId"`
	OwnerPartyId    uint32 `json:"ownerPartyId"`
	DropTime        uint64 `json:"dropTime"`
	DropperUniqueId uint32 `json:"dropperUniqueId"`
	DropperX        int16  `json:"dropperX"`
	DropperY        int16  `json:"dropperY"`
	PlayerDrop      bool   `json:"playerDrop"`
	Mod             byte   `json:"mod"`
}

func DropEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &DropEvent{}
	}
}

func HandleDropEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*DropEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			session.ForEachInMap(wid, cid, event.MapId, dropItem(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func dropItem(l logrus.FieldLogger, event *DropEvent) session.Operator {
	return func(s *session.Model) {
		a := uint32(0)
		if event.ItemId != 0 {
			a = 0
		} else {
			a = event.Meso
		}
		err := s.Announce(writer.WriteDropItemFromMapObject(l)(event.UniqueId, event.ItemId, event.Meso, a,
			event.DropperUniqueId, event.DropType, event.OwnerId, event.OwnerPartyId, s.CharacterId(), 0,
			event.DropTime, event.DropX, event.DropY, event.DropperX, event.DropperY, event.PlayerDrop, event.Mod))
		if err != nil {
			l.WithError(err).Errorf("Unable to write drop in map for character %d", s.CharacterId())
		}
	}
}
