package writer

import "atlas-wcc/socket/response"

const OpCodeShowForeignEffect uint16 = 0xC6

func WriteShowForeignEffect(characterId uint32, effect byte) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeShowForeignEffect)
	w.WriteInt(characterId)
	w.WriteByte(effect)
	return w.Bytes()
}
