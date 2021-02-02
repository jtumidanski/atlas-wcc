package writer

import "atlas-wcc/socket/response"

const OpCodeSetWeekEventMessage uint16 = 0x4D

func WriteYellowTip(message string) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeSetWeekEventMessage)
	w.WriteByte(0xFF)
	w.WriteAsciiString(message)
	w.WriteShort(0)
	return w.Bytes()
}
