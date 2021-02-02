package writer

import "atlas-wcc/socket/response"

const OpCodeChatText uint16 = 0xA2

func WriteChatText(characterId uint32, message string, gm bool, show bool) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeChatText)
	w.WriteInt(characterId)
	w.WriteBool(gm)
	w.WriteAsciiString(message)
	w.WriteBool(show)
	return w.Bytes()
}
