package keymap

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeKeyMap uint16 = 0x14F

func WriteKeyMap(l logrus.FieldLogger) func(keys []*Model) []byte {
	return func(keys []*Model) []byte {
		km := make(map[int32]*Model)
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
