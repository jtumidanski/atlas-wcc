package server

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeServerMessage uint16 = 0x44
const OpCodeSetWeekEventMessage uint16 = 0x4D

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

func WritePinkText(l logrus.FieldLogger) func(message string) []byte {
	return func(message string) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeServerMessage)
		w.WriteByte(5)
		w.WriteAsciiString(message)
		return w.Bytes()
	}
}

func WritePopup(l logrus.FieldLogger) func(message string) []byte {
	return func(message string) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeServerMessage)
		w.WriteByte(1)
		w.WriteAsciiString(message)
		return w.Bytes()
	}
}

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
