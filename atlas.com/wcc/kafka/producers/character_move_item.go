package producers

import "github.com/sirupsen/logrus"

type moveItemCommand struct {
	CharacterId   uint32 `json:"characterId"`
	InventoryType int8   `json:"inventoryType"`
	Source        int16  `json:"source"`
	Destination   int16  `json:"destination"`
}

func MoveItem(l logrus.FieldLogger) func(characterId uint32, inventoryType int8, source int16, destination int16) {
	producer := ProduceEvent(l, "TOPIC_MOVE_ITEM")
	return func(characterId uint32, inventoryType int8, source int16, destination int16) {
		e := &moveItemCommand{CharacterId: characterId, InventoryType: inventoryType, Source: source, Destination: destination}
		producer(CreateKey(int(characterId)), e)
	}
}
