package producers

import (
	"github.com/sirupsen/logrus"
)

type characterDistributeSpEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
}

func CharacterDistributeSp(l logrus.FieldLogger) func(characterId uint32, skillId uint32) {
	producer := ProduceEvent(l, "TOPIC_ASSIGN_SP_COMMAND")
	return func(characterId uint32, skillId uint32) {
		e := &characterDistributeSpEvent{
			CharacterId: characterId,
			SkillId:     skillId,
		}
		producer(CreateKey(int(characterId)), e)
	}
}
