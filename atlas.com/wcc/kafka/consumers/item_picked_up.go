package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*itemPickedUpEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, showItemGain(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showItemGain(_ logrus.FieldLogger, event *itemPickedUpEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteShowItemGain(event.ItemId, event.Quantity))
		session.Announce(writer.WriteEnableActions())
	}
}
