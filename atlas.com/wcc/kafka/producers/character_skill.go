package producers

import "github.com/sirupsen/logrus"

type applySkillCommand struct {
	CharacterId uint32
	SkillId     uint32
	Level       uint8
	X           int16
	Y           int16
}

func ApplySkill(l logrus.FieldLogger) func(characterId uint32, skillId uint32, level uint8, x int16, y int16) {
	producer := ProduceEvent(l, "TOPIC_APPLY_SKILL_COMMAND")
	return func(characterId uint32, skillId uint32, level uint8, x int16, y int16) {
		e := &applySkillCommand{
			CharacterId: characterId,
			SkillId:     skillId,
			Level:       level,
			X:           x,
			Y:           y,
		}
		producer(CreateKey(int(characterId)), e)
	}
}

type MonsterMagnetData struct {
	MonsterId uint32
	Success   uint8
}

type applyMonsterMagnetCommand struct {
	CharacterId uint32
	SkillId     uint32
	Level       uint8
	Direction   int8
	Data        []MonsterMagnetData
}

func ApplyMonsterMagnet(l logrus.FieldLogger) func(characterId uint32, skillId uint32, level uint8, direction int8, data []MonsterMagnetData) {
	producer := ProduceEvent(l, "TOPIC_APPLY_MONSTER_MAGNET_COMMAND")
	return func(characterId uint32, skillId uint32, level uint8, direction int8, data []MonsterMagnetData) {
		e := applyMonsterMagnetCommand{
			CharacterId: characterId,
			SkillId:     skillId,
			Level:       level,
			Direction:   direction,
			Data:        data,
		}
		producer(CreateKey(int(characterId)), e)
	}
}
