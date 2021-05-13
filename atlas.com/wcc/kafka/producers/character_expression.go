package producers

import (
	"github.com/sirupsen/logrus"
)

type characterExpressionEvent struct {
	CharacterId uint32 `json:"characterId"`
	Emote       uint32 `json:"emote"`
}

func CharacterExpression(l logrus.FieldLogger) func(characterId uint32, emote uint32) {
	producer := ProduceEvent(l, "CHANGE_FACIAL_EXPRESSION")
	return func(characterId uint32, emote uint32) {
		e := &characterExpressionEvent{
			CharacterId: characterId,
			Emote:       emote,
		}
		producer(CreateKey(int(characterId)), e)
	}
}