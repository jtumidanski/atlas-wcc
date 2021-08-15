package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeReactorTrigger uint16 = 0x115

func WriteReactorTrigger(l logrus.FieldLogger) func(id uint32, state int8, x int16, y int16, stance byte) []byte {
	return func(id uint32, state int8, x int16, y int16, stance byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeReactorTrigger)
		w.WriteInt(id)
		w.WriteInt8(state)
		w.WriteInt16(x)
		w.WriteInt16(y)
		w.WriteByte(stance)
		w.WriteByte(0)
		w.WriteInt16(5) // frame delay, set to 5 since there doesn't appear to be a fixed formula for it
		return w.Bytes()
	}
}