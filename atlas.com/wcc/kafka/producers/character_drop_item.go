package producers

import "github.com/sirupsen/logrus"

type characterDropItem struct {
	WorldId       byte   `json:"worldId"`
	ChannelId     byte   `json:"channelId"`
	CharacterId   uint32 `json:"characterId"`
	InventoryType int8   `json:"inventoryType"`
	Source        int16  `json:"source"`
	Quantity      int16  `json:"quantity"`
}

func DropItem(l logrus.FieldLogger) func(worldId byte, channelId byte, chcharacterId uint32, inventoryType int8, source int16, quantity int16) {
	return func(worldId byte, channelId byte, characterId uint32, inventoryType int8, source int16, quantity int16) {
		e := &characterDropItem{WorldId: worldId, ChannelId: channelId, CharacterId: characterId, InventoryType: inventoryType, Source: source, Quantity: quantity}
		produceEvent(l, "TOPIC_CHARACTER_DROP_ITEM", createKey(int(characterId)), e)
	}
}
