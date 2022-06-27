package server

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeServerMessage uint16 = 0x44

func WriteServerNotice(l logrus.FieldLogger) func(channelId byte, t byte, m string, mega bool, npc uint32) []byte {
	return func(channelId byte, t byte, m string, mega bool, npc uint32) []byte {
		w := response.NewWriter(l)
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
}
