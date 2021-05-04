package writer

import "atlas-wcc/socket/response"

const OpCodeCancelBuff uint16 = 0x21

func WriteCancelBuff(stats []BuffStat) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeCancelBuff)
	writeLongMask(w, stats)
	w.WriteByte(1)
	return w.Bytes()
}