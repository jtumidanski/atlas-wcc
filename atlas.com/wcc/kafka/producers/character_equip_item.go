package producers

import "github.com/sirupsen/logrus"

type characterEquipItem struct {
	CharacterId uint32 `json:"characterId"`
	Source      int16  `json:"source"`
	Destination int16  `json:"destination"`
}

func EquipItem(l logrus.FieldLogger) func(characterId uint32, source int16, destination int16) {
	producer := ProduceEvent(l, "TOPIC_EQUIP_ITEM")
	return func(characterId uint32, source int16, destination int16) {
		e := &characterEquipItem{CharacterId: characterId, Source: source, Destination: destination}
		producer(CreateKey(int(characterId)), e)
	}
}
