package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
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

func DropEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &DropEvent{}
	}
}

func HandleDropEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*DropEvent); ok {
			sl, err := getSessionsForThoseInMap(wid, cid, event.MapId)
			if err != nil {
				return
			}
			a := uint32(0)
			if event.ItemId != 0 {
				a = 0
			} else {
				a = event.Meso
			}
			for _, s := range sl {
				s.Announce(writer.WriteDropItemFromMapObject(event.UniqueId, event.ItemId, event.Meso, a, event.DropperUniqueId, event.DropType, event.OwnerId, event.OwnerPartyId, s.CharacterId(), 0, event.DropTime, event.DropX, event.DropY, event.DropperX, event.DropperY, event.PlayerDrop, event.Mod))
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleDropEvent]")
		}
	}
}