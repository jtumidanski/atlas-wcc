package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeShowForeignEffect uint16 = 0xC6

func WriteShowForeignEffect(l logrus.FieldLogger) func(characterId uint32, effect byte) []byte {
	return func(characterId uint32, effect byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowForeignEffect)
		w.WriteInt(characterId)
		w.WriteByte(effect)
		return w.Bytes()
	}
}

func WriteShowBuffEffect(l logrus.FieldLogger) func(characterId uint32, effect byte, skillId uint32, direction byte) []byte {
	return func(characterId uint32, effect byte, skillId uint32, direction byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowForeignEffect)
		w.WriteInt(characterId)
		w.WriteByte(effect)
		w.WriteInt(skillId)
		w.WriteByte(direction)
		return w.Bytes()
	}
}
