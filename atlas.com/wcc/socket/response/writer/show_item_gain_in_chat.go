package writer

import "atlas-wcc/socket/response"

const OpCodeShowItemGainInChat uint16 = 0xCE

func WriteShowOwnBuff(effect byte, skillId uint32) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeShowItemGainInChat)
	w.WriteByte(effect)
	w.WriteInt(skillId)
	w.WriteByte(0xA9)
	w.WriteByte(1)
	return w.Bytes()
}
