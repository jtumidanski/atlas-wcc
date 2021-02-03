package writer

import "atlas-wcc/socket/response"

const OpCodeCharacterUpdateSkills uint16 = 0x24

func WriteCharacterSkillUpdate(skillId uint32, skillLevel uint32, masterLevel uint32, expiration uint64) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeCharacterUpdateSkills)
	w.WriteByte(1)
	w.WriteShort(1)
	w.WriteInt(skillId)
	w.WriteInt(skillLevel)
	w.WriteInt(masterLevel)
	w.WriteLong(expiration)
	w.WriteByte(4)
	return w.Bytes()
}
