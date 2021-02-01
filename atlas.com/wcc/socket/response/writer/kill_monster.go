package writer

import (
	"atlas-wcc/socket/response"
)

const OpCodeKillMonster uint16 = 0xED

func WriteKillMonster(uniqueId uint32, animation bool) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeKillMonster)
	w.WriteInt(uniqueId)
	w.WriteBool(animation)
	w.WriteBool(animation)
	return w.Bytes()
}
