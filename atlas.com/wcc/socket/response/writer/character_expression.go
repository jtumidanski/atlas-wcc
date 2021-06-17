package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeCharacterExpression uint16 = 0xC1

func WriteCharacterExpression(l logrus.FieldLogger) func(characterId uint32, expression uint32) []byte {
	return func(characterId uint32, expression uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCharacterExpression)
		w.WriteInt(characterId)
		w.WriteInt(expression)
		return w.Bytes()
	}
}
