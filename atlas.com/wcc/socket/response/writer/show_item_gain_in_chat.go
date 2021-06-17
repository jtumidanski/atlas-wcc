package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeShowItemGainInChat uint16 = 0xCE

func WriteShowOwnBuff(l logrus.FieldLogger) func(effect byte, skillId uint32) []byte {
	return func(effect byte, skillId uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowItemGainInChat)
		w.WriteByte(effect)
		w.WriteInt(skillId)
		w.WriteByte(0xA9)
		w.WriteByte(1)
		return w.Bytes()
	}
}
