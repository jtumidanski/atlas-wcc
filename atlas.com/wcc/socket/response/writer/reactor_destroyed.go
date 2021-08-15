package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeReactorDestroyed uint16 = 0x118

func WriteReactorDestroyed(l logrus.FieldLogger) func(id uint32, state int8, x int16, y int16) []byte {
	return func(id uint32, state int8, x int16, y int16) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeReactorDestroyed)
		w.WriteInt(id)
		w.WriteInt8(state)
		w.WriteInt16(x)
		w.WriteInt16(y)
		return w.Bytes()
	}
}
