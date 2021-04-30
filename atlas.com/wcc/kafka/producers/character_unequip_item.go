package producers

import "github.com/sirupsen/logrus"

type characterUnequipItem struct {
	CharacterId uint32 `json:"characterId"`
	Source      int16  `json:"source"`
	Destination int16  `json:"destination"`
}

func UnequipItem(l logrus.FieldLogger) func(characterId uint32, source int16, destination int16) {
	return func(characterId uint32, source int16, destination int16) {
		e := &characterUnequipItem{CharacterId: characterId, Source: source, Destination: destination}
		produceEvent(l, "TOPIC_UNEQUIP_ITEM", createKey(int(characterId)), e)
	}
}
