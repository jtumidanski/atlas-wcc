package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if command, ok := e.(*inventoryFullCommand); ok {
			session.ForSessionByCharacterId(command.CharacterId, showInventoryFull(l))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showInventoryFull(l logrus.FieldLogger) session.SessionOperator {
	return func(s *session.Model) {
		err := s.Announce(writer.WriteShowInventoryFull())
		if err != nil {
			l.WithError(err).Errorf("Unable to show inventory is full for character %d.", s.CharacterId())
		}
	}
}
