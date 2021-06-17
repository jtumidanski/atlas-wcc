package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeKillMonster uint16 = 0xED

func WriteKillMonster(l logrus.FieldLogger) func(uniqueId uint32, animation bool) []byte {
	return func(uniqueId uint32, animation bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeKillMonster)
		w.WriteInt(uniqueId)
		w.WriteBool(animation)
		w.WriteBool(animation)
		return w.Bytes()
	}
}
