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
			if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForEachInMap(l)(wid, cid, event.MapId, removeItem(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func removeItem(l logrus.FieldLogger, event *dropPickedUpEvent) session.Operator {
	b := writer.WriteRemoveItem(l)(event.DropId, 2, event.CharacterId)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
