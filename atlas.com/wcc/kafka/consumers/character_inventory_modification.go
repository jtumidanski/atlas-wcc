package consumers

import (
	"atlas-wcc/domain"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type inventoryModification struct {
	Mode          byte   `json:"mode"`
	ItemId        uint32 `json:"itemId"`
	InventoryType byte   `json:"inventoryType"`
	Quantity      uint16 `json:"quantity"`
	Position      int16  `json:"position"`
	OldPosition   int16  `json:"oldPosition"`
}

type characterInventoryModificationEvent struct {
	CharacterId   uint32                  `json:"characterId"`
	UpdateTick    bool                    `json:"updateTick"`
	Modifications []inventoryModification `json:"modifications"`
}

func CharacterInventoryModificationEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterInventoryModificationEvent{}
	}
}

func HandleCharacterInventoryModificationEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterInventoryModificationEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, writeInventoryModification(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeInventoryModification(l logrus.FieldLogger, event *characterInventoryModificationEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		result := writer.ModifyInventory{}
		result.UpdateTick = event.UpdateTick
		for _, m := range event.Modifications {
			var item writer.InventoryItem
			if m.InventoryType == 1 {
				e, err := processors.GetEquipItemForCharacter(event.CharacterId, m.Position)
				if err != nil {
					l.WithError(err).Errorf("Retrieving equipment in position %d for character %d.", m.Position, event.CharacterId)
					continue
				}
				item = e
			} else {
				item = domain.NewItem(m.ItemId, m.Position, m.Quantity)
			}

			mi := writer.Modification{
				Mode:          m.Mode,
				InventoryType: m.InventoryType,
				Item:          item,
				OldPosition:   m.OldPosition,
			}
			result.Modifications = append(result.Modifications, mi)
		}
		session.Announce(writer.WriteCharacterInventoryModification(result))
	}
}
