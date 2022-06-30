package inventory

import (
	"atlas-wcc/kafka"
	"atlas-wcc/model"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameModification = "inventory_modification_event"
	consumerNameFull         = "inventory_full_command"
	topicTokenModification   = "TOPIC_INVENTORY_MODIFICATION"
	topicTokenFull           = "TOPIC_INVENTORY_FULL"
)

func ModificationConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[modificationEvent](consumerNameModification, topicTokenModification, groupId, handleModification(wid, cid))
	}
}

type modificationEvent struct {
	CharacterId   uint32         `json:"characterId"`
	UpdateTick    bool           `json:"updateTick"`
	Modifications []modification `json:"modifications"`
}

type modification struct {
	Mode          byte   `json:"mode"`
	ItemId        uint32 `json:"itemId"`
	InventoryType int8   `json:"inventoryType"`
	Quantity      uint16 `json:"quantity"`
	Position      int16  `json:"position"`
	OldPosition   int16  `json:"oldPosition"`
}

func handleModification(_ byte, _ byte) kafka.HandlerFunc[modificationEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event modificationEvent) {
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}
		session.IfPresentByCharacterId(event.CharacterId, writeModification(l, span)(event))
	}
}

func writeModification(l logrus.FieldLogger, span opentracing.Span) func(event modificationEvent) model.Operator[session.Model] {
	return func(event modificationEvent) model.Operator[session.Model] {
		result := ModifyInventory{}
		result.UpdateTick = event.UpdateTick
		for _, m := range event.Modifications {
			var item InventoryItem
			if m.InventoryType == 1 {
				if m.Mode == 3 {
					// create dummy item for removal.
					item = NewItem(m.ItemId, m.Position, 1)
				} else {
					e, err := GetEquipItemForCharacter(l, span)(event.CharacterId, m.Position)
					if err != nil {
						l.WithError(err).Errorf("Retrieving equipment in position %d for character %d.", m.Position, event.CharacterId)
						continue
					}
					item = e
				}
			} else {
				item = NewItem(m.ItemId, m.Position, m.Quantity)
			}

			mi := Modification{
				Mode:          m.Mode,
				InventoryType: m.InventoryType,
				Item:          item,
				OldPosition:   m.OldPosition,
			}
			result.Modifications = append(result.Modifications, mi)
		}
		return session.AnnounceOperator(WriteCharacterInventoryModification(l)(result))
	}
}

func FullConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[fullCommand](consumerNameFull, topicTokenFull, groupId, handleFull(wid, cid))
	}
}

type fullCommand struct {
	CharacterId uint32 `json:"characterId"`
}

func handleFull(_ byte, _ byte) kafka.HandlerFunc[fullCommand] {
	return func(l logrus.FieldLogger, span opentracing.Span, command fullCommand) {
		session.IfPresentByCharacterId(command.CharacterId, session.AnnounceOperator(WriteShowInventoryFull(l)))
	}
}
