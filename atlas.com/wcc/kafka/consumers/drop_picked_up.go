package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type dropPickedUpEvent struct {
	DropId      uint32 `json:"dropId"`
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
}

func DropPickedUpEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &dropPickedUpEvent{}
	}
}

func HandleDropPickedUpEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*dropPickedUpEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForEachSessionInMap(wid, cid, event.MapId, removeItem(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func removeItem(_ logrus.FieldLogger, event *dropPickedUpEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteRemoveItem(event.DropId, 2, event.CharacterId))
	}
}
