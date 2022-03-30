package skill

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeCharacterUpdateSkills uint16 = 0x24

func WriteCharacterSkillUpdate(l logrus.FieldLogger) func(skillId uint32, skillLevel uint32, masterLevel uint32, expiration int64) []byte {
	return func(skillId uint32, skillLevel uint32, masterLevel uint32, expiration int64) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCharacterUpdateSkills)
		w.WriteByte(1)
		w.WriteShort(1)
		w.WriteInt(skillId)
		w.WriteInt(skillLevel)
		w.WriteInt(masterLevel)
		w.WriteInt64(expiration)
		w.WriteByte(4)
		return w.Bytes()
	}
}
