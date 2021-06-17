package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeCancelBuff uint16 = 0x21

func WriteCancelBuff(l logrus.FieldLogger) func(stats []BuffStat) []byte {
	return func(stats []BuffStat) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCancelBuff)
		writeLongMask(w, stats)
		w.WriteByte(1)
		return w.Bytes()
	}
}
