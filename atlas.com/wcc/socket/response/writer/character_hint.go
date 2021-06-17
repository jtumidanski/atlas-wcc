package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeCharacterHint uint16 = 0xD6

func WriteHint(l logrus.FieldLogger) func(message string, width int16, height int16) []byte {
	return func(message string, width int16, height int16) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCharacterHint)
		tw := width
		th := height
		if tw < 1 {
			tw = int16(len(message) * 10)
			if tw < 40 {
				tw = 40
			}
		}
		if th < 5 {
			th = 5
		}

		w.WriteAsciiString(message)
		w.WriteInt16(tw)
		w.WriteInt16(th)
		w.WriteByte(1)
		return w.Bytes()
	}
}
