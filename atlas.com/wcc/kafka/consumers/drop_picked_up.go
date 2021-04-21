package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*dropPickedUpEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachSessionInMap(wid, cid, event.MapId, removeItem(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func removeItem(_ logrus.FieldLogger, event *dropPickedUpEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteRemoveItem(event.DropId, 2, event.CharacterId))
	}
}
