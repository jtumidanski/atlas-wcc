package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type itemPickedUpEvent struct {
	CharacterId uint32 `json:"characterId"`
	ItemId      uint32 `json:"itemId"`
	Quantity    uint32 `json:"quantity"`
}

func ItemPickedUpEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &itemPickedUpEvent{}
	}
}

func HandleItemPickedUpEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*itemPickedUpEvent); ok {
			if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, showItemGain(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showItemGain(l logrus.FieldLogger, event *itemPickedUpEvent) session.Operator {
	ig := writer.WriteShowItemGain(l)(event.ItemId, event.Quantity)
	ea := writer.WriteEnableActions(l)
	return func(s *session.Model) {
		err := s.Announce(ig)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		err = s.Announce(ea)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
