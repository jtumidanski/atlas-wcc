package writer

import "atlas-wcc/socket/response"

const OpCodeServerMessage uint16 = 0x44

func WriteServerNotice(channelId byte, t byte, m string, mega bool, npc uint32) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeServerMessage)
	w.WriteByte(t)
	w.WriteAsciiString(m)
	if t == 3 {
		w.WriteByte(channelId - 1)
		w.WriteBool(mega)
	} else if t == 6 {
		w.WriteInt(0)
	} else if t == 7 {
		w.WriteInt(npc)
	}
	return w.Bytes()
}
