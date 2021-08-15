package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeReactorSpawn uint16 = 0x117

func WriteReactorSpawn(l logrus.FieldLogger) func(id uint32, classification uint32, state int8, x int16, y int16) []byte {
	return func(id uint32, classification uint32, state int8, x int16, y int16) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeReactorSpawn)
		w.WriteInt(id)
		w.WriteInt(classification)
		w.WriteInt8(state)
		w.WriteInt16(x)
		w.WriteInt16(y)
		w.WriteByte(0)
		w.WriteInt16(0)
		return w.Bytes()
	}
}
