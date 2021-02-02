package producers

import (
	"context"
	"log"
)

type MonsterMovement struct {
	l   *log.Logger
	ctx context.Context
}

func NewMonsterMovement(l *log.Logger, ctx context.Context) *MonsterMovement {
	return &MonsterMovement{l, ctx}
}

func (m *MonsterMovement) EmitMovement(uniqueId uint32, observerId uint32, skillPossible bool, skill int8, skillId uint32, skillLevel uint32, option uint16, startX int16, startY int16, endX int16, endY int16, stance byte, rawMovement []byte) {
	e := &MonsterMovementEvent{
		UniqueId:      uniqueId,
		ObserverId:    observerId,
		SkillPossible: skillPossible,
		Skill:         skill,
		SkillId:       skillId,
		SkillLevel:    skillLevel,
		Option:        option,
		StartX:        startX,
		StartY:        startY,
		EndX:          endX,
		EndY:          endY,
		Stance:        stance,
		RawMovement:   rawMovement,
	}
	ProduceEvent(m.l, "TOPIC_MONSTER_MOVEMENT", createKey(int(uniqueId)), e)
}

type MonsterMovementEvent struct {
	UniqueId      uint32      `json:"uniqueId"`
	ObserverId    uint32      `json:"observerId"`
	SkillPossible bool        `json:"skillPossible"`
	Skill         int8        `json:"skill"`
	SkillId       uint32      `json:"skillId"`
	SkillLevel    uint32      `json:"skillLevel"`
	Option        uint16      `json:"option"`
	StartX        int16       `json:"startX"`
	StartY        int16       `json:"startY"`
	EndX          int16       `json:"endX"`
	EndY          int16       `json:"endY"`
	Stance        byte        `json:"stance"`
	RawMovement   RawMovement `json:"rawMovement"`
}
