package writer

import "atlas-wcc/socket/response"

const OpCodeMoveMonsterResponse uint16 = 0xF0

func WriteMoveMonsterResponse(objectId uint32, moveId uint16, currentMp uint16, useSkills bool, skillId byte, skillLevel byte) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeMoveMonsterResponse)
	w.WriteInt(objectId)
	w.WriteShort(moveId)
	w.WriteBool(useSkills)
	w.WriteShort(currentMp)
	w.WriteByte(skillId)
	w.WriteByte(skillLevel)
	return w.Bytes()
}