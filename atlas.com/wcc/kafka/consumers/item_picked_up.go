package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type itemPickedUpEvent struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    uint32 `json:"quantity"`
}

func ItemPickedUpEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &itemPickedUpEvent{}
	}
}

func HandleItemPickedUpEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*itemPickedUpEvent); ok {
			processors.ForSessionByCharacterId(event.CharacterId, showItemGain(l, event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleItemPickedUpEvent]")
		}
	}
}

func showItemGain(_ *log.Logger, event *itemPickedUpEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteShowItemGain(event.ItemId, event.Quantity))
		session.Announce(writer.WriteEnableActions())
	}
}
