package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeNPCTalk uint16 = 0x130

func WriteNPCTalk(l logrus.FieldLogger) func(npcId uint32, messageType byte, talk string, end []byte, speaker byte) []byte {
	return func(npcId uint32, messageType byte, talk string, end []byte, speaker byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCTalk)
		w.WriteByte(4)
		w.WriteInt(npcId)
		w.WriteByte(messageType)
		w.WriteByte(speaker)
		w.WriteAsciiString(talk)
		w.WriteByteArray(end)
		return w.Bytes()
	}
}

func WriteNPCTalkNum(l logrus.FieldLogger) func(npcId uint32, talk string, defaultValue int32, minimumValue int32, maximumValue int32) []byte {
	return func(npcId uint32, talk string, defaultValue int32, minimumValue int32, maximumValue int32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCTalk)
		w.WriteByte(4)
		w.WriteInt(npcId)
		w.WriteByte(3)
		w.WriteByte(0)
		w.WriteAsciiString(talk)
		w.WriteInt32(defaultValue)
		w.WriteInt32(minimumValue)
		w.WriteInt32(maximumValue)
		w.WriteInt(0)
		return w.Bytes()
	}
}

func WriteNPCTalkStyle(l logrus.FieldLogger) func(npcId uint32, talk string, styles []uint32) []byte {
	return func(npcId uint32, talk string, styles []uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCTalk)
		w.WriteByte(4)
		w.WriteInt(npcId)
		w.WriteByte(7)
		w.WriteByte(0)
		w.WriteAsciiString(talk)
		w.WriteByte(byte(len(styles)))
		for _, style := range styles {
			w.WriteInt(style)
		}
		return w.Bytes()
	}
}
