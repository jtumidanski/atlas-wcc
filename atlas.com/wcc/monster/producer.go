package monster

import (
	"atlas-wcc/kafka"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strings"
)

type movementEvent struct {
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
	RawMovement   rawMovement `json:"rawMovement"`
}

type rawMovement []byte

func (m rawMovement) MarshalJSON() ([]byte, error) {
	var result string
	if m == nil {
		result = "[]"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", m)), ",")
	}
	return []byte(result), nil
}

func Move(l logrus.FieldLogger, span opentracing.Span) func(uniqueId uint32, observerId uint32, skillPossible bool, skill int8, skillId uint32, skillLevel uint32, option uint16, startX int16, startY int16, endX int16, endY int16, stance byte, rawMovement []byte) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_MONSTER_MOVEMENT")
	return func(uniqueId uint32, observerId uint32, skillPossible bool, skill int8, skillId uint32, skillLevel uint32, option uint16, startX int16, startY int16, endX int16, endY int16, stance byte, rawMovement []byte) {
		e := &movementEvent{
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
		producer(kafka.CreateKey(int(uniqueId)), e)
	}
}
