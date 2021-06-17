package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeChatText uint16 = 0xA2

func WriteChatText(l logrus.FieldLogger) func(characterId uint32, message string, gm bool, show bool) []byte {
	return func(characterId uint32, message string, gm bool, show bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeChatText)
		w.WriteInt(characterId)
		w.WriteBool(gm)
		w.WriteAsciiString(message)
		w.WriteBool(show)
		return w.Bytes()
	}
}
