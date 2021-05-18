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

func WriteShowBuffEffect(characterId uint32, effect byte, skillId uint32, direction byte) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeShowForeignEffect)
	w.WriteInt(characterId)
	w.WriteByte(effect)
	w.WriteInt(skillId)
	w.WriteByte(direction)
	return w.Bytes()
}
