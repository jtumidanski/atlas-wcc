package producers

import "github.com/sirupsen/logrus"

type applySkillCommand struct {
	characterId uint32
	skillId     uint32
	level       uint8
	x           int16
	y           int16
}

func ApplySkill(l logrus.FieldLogger) func(characterId uint32, skillId uint32, level uint8, x int16, y int16) {
	return func(characterId uint32, skillId uint32, level uint8, x int16, y int16) {
		e := &applySkillCommand{
			characterId: characterId,
			skillId:     skillId,
			level:       level,
			x:           x,
			y:           y,
		}
		produceEvent(l, "TOPIC_APPLY_SKILL_COMMAND", createKey(int(characterId)), e)
	}
}

type MonsterMagnetData struct {
	MonsterId uint32
	Success   uint8
}

type applyMonsterMagnetCommand struct {
	characterId uint32
	skillId     uint32
	level       uint8
	direction   int8
	data        []MonsterMagnetData
}

func ApplyMonsterMagnet(l logrus.FieldLogger) func(characterId uint32, skillId uint32, level uint8, direction int8, data []MonsterMagnetData) {
	return func(characterId uint32, skillId uint32, level uint8, direction int8, data []MonsterMagnetData) {
		e := applyMonsterMagnetCommand{
			characterId: characterId,
			skillId:     skillId,
			level:       level,
			direction:   direction,
			data:        data,
		}
		produceEvent(l, "TOPIC_APPLY_MONSTER_MAGNET_COMMAND", createKey(int(characterId)), e)
	}
}
