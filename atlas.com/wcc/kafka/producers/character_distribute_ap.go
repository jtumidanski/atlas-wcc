package producers

import (
	"github.com/sirupsen/logrus"
)

type characterDistributeApEvent struct {
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func CharacterDistributeAp(l logrus.FieldLogger) func(characterId uint32, attributeType string) {
	producer := ProduceEvent(l, "TOPIC_ASSIGN_AP_COMMAND")
	return func(characterId uint32, attributeType string) {
		e := &characterDistributeApEvent{
			CharacterId: characterId,
			Type:        attributeType,
		}
		producer(CreateKey(int(characterId)), e)
	}
}
