package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeRemoveItemFromMap uint16 = 0x10D

func WriteRemoveItem(l logrus.FieldLogger) func(dropId uint32, animation byte, characterId uint32) []byte {
	return func(dropId uint32, animation byte, characterId uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeRemoveItemFromMap)
		w.WriteByte(animation)
		w.WriteInt(dropId)
		if animation >= 2 {
			w.WriteInt(characterId)
			//TODO handle pet loot
		}
		return w.Bytes()
	}
}
