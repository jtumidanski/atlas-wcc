package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type inventoryFullCommand struct {
	CharacterId uint32 `json:"characterId"`
}

func InventoryFullCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &inventoryFullCommand{}
	}
}

func HandleInventoryFullCommand() ChannelEventProcessor {
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if command, ok := e.(*inventoryFullCommand); ok {
			session.ForSessionByCharacterId(command.CharacterId, showInventoryFull(l))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showInventoryFull(l logrus.FieldLogger) session.Operator {
	return func(s *session.Model) {
		err := s.Announce(writer.WriteShowInventoryFull(l))
		if err != nil {
			l.WithError(err).Errorf("Unable to show inventory is full for character %d.", s.CharacterId())
		}
	}
}
