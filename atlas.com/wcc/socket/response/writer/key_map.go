package writer

import (
	"atlas-wcc/character/keymap"
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeKeyMap uint16 = 0x14F

func WriteKeyMap(l logrus.FieldLogger) func(keys []*keymap.Model) []byte {
	return func(keys []*keymap.Model) []byte {
		km := make(map[int32]*keymap.Model)
		for _, k := range keys {
			km[k.Key()] = k
		}

		w := response.NewWriter(l)
		w.WriteShort(OpCodeKeyMap)
		w.WriteByte(0)
		for i := int32(0); i < 90; i++ {
			if k, ok := km[i]; ok {
				w.WriteInt8(k.Type())
				w.WriteInt32(k.Action())
			} else {
				w.WriteInt8(0)
				w.WriteInt32(0)
			}
		}
		return w.Bytes()
	}
}
