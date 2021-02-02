package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
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
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterInventoryModificationEvent); ok {
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}

			result := writer.ModifyInventory{}
			result.UpdateTick = event.UpdateTick
			for _, m := range event.Modifications {
				mi := writer.Modification{
					Mode:          m.Mode,
					ItemId:        m.ItemId,
					InventoryType: m.InventoryType,
					Quantity:      m.Quantity,
					Position:      m.Position,
					OldPosition:   m.OldPosition,
				}
				result.Modifications = append(result.Modifications, mi)
			}

			(*as).Announce(writer.WriteCharacterInventoryModification(result))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterInventoryModificationEvent]")
		}
	}
}
