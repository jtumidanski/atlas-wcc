package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeMoveMonsterResponse uint16 = 0xF0

func WriteMoveMonsterResponse(l logrus.FieldLogger) func(objectId uint32, moveId uint16, currentMp uint16, useSkills bool, skillId byte, skillLevel byte) []byte {
	return func(objectId uint32, moveId uint16, currentMp uint16, useSkills bool, skillId byte, skillLevel byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeMoveMonsterResponse)
		w.WriteInt(objectId)
		w.WriteShort(moveId)
		w.WriteBool(useSkills)
		w.WriteShort(currentMp)
		w.WriteByte(skillId)
		w.WriteByte(skillLevel)
		return w.Bytes()
	}
}
