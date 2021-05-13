package producers

import (
	"github.com/sirupsen/logrus"
)

type characterAdjustManaEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

func CharacterAdjustMana(l logrus.FieldLogger) func(characterId uint32, amount int16) {
	producer := ProduceEvent( l, "TOPIC_ADJUST_MANA")
	return func(characterId uint32, amount int16) {
		e := &characterAdjustManaEvent{
			CharacterId: characterId,
			Amount:      amount,
		}
		producer(CreateKey(int(characterId)), e)
	}
}