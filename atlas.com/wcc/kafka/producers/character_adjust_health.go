package producers

import (
	"github.com/sirupsen/logrus"
)

type characterAdjustHealthEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

func CharacterAdjustHealth(l logrus.FieldLogger) func(characterId uint32, amount int16) {
	producer := ProduceEvent(l, "TOPIC_ADJUST_HEALTH")
	return func(characterId uint32, amount int16) {
		e := &characterAdjustHealthEvent{
			CharacterId: characterId,
			Amount:      amount,
		}
		producer(CreateKey(int(characterId)), e)
	}
}
