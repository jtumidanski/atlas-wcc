package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeNPCAction uint16 = 0x104

func WriteNPCAnimation(l logrus.FieldLogger) func(objectId uint32, second byte, third byte) []byte {
	return func(objectId uint32, second byte, third byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCAction)
		w.WriteInt(objectId)
		w.WriteByte(second)
		w.WriteByte(third)
		return w.Bytes()
	}
}

func WriteNPCMove(l logrus.FieldLogger) func(movement []byte) []byte {
	return func(movement []byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCAction)
		w.WriteByteArray(movement)
		return w.Bytes()
	}
}
