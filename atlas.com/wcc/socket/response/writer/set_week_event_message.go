package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeSetWeekEventMessage uint16 = 0x4D

func WriteYellowTip(l logrus.FieldLogger) func(message string) []byte {
	return func(message string) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSetWeekEventMessage)
		w.WriteByte(0xFF)
		w.WriteAsciiString(message)
		w.WriteShort(0)
		return w.Bytes()
	}
}
