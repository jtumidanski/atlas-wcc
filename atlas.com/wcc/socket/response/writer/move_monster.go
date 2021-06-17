package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeMoveMonster uint16 = 0xEF

func WriteMoveMonster(l logrus.FieldLogger) func(objectId uint32, skillPossible bool, skill int8, skillId byte, skillLevel byte, option uint16, startX int16, startY int16, movementList []byte) []byte {
	return func(objectId uint32, skillPossible bool, skill int8, skillId byte, skillLevel byte, option uint16, startX int16, startY int16, movementList []byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeMoveMonster)
		w.WriteInt(objectId)
		w.WriteByte(0)
		w.WriteBool(skillPossible)
		w.WriteInt8(skill)
		w.WriteByte(skillId)
		w.WriteByte(skillLevel)
		w.WriteShort(option)
		w.WriteInt16(startX)
		w.WriteInt16(startY)
		for _, b := range movementList {
			w.WriteByte(b)
		}
		return w.Bytes()
	}
}
