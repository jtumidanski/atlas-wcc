package writer

import "atlas-wcc/socket/response"

const OpCodeNPCTalk uint16 = 0x130

func WriteNPCTalk(npcId uint32, messageType byte, talk string, end []byte, speaker byte) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeNPCTalk)
	w.WriteByte(4)
	w.WriteInt(npcId)
	w.WriteByte(messageType)
	w.WriteByte(speaker)
	w.WriteAsciiString(talk)
	w.WriteByteArray(end)
	return w.Bytes()
}
