package writer

import "atlas-wcc/socket/response"

const OpCodeCharacterDamage uint16 = 0xC0

func WriteCharacterDamaged(skillId int8, monsterId uint32, characterId uint32, damage uint32, fake uint32, direction int8, pgmr bool, pgmr1 byte, pg bool, monsterUniqueId uint32, x int16, y int16) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeCharacterDamage)
	w.WriteInt(characterId)
	w.WriteInt8(skillId)
	w.WriteInt(damage)
	if skillId != -4 {
		w.WriteInt(monsterId)
		w.WriteInt8(direction)
		if pgmr {
			w.WriteByte(pgmr1)
			if pg {
				w.WriteByte(1)
			} else {
				w.WriteByte(0)
			}
			w.WriteInt(monsterUniqueId)
			w.WriteByte(6)
			w.WriteInt16(x)
			w.WriteInt16(y)
			w.WriteByte(0)
		} else {
			w.WriteShort(0)
		}
		w.WriteInt(damage)
		if fake > 0 {
			w.WriteInt(fake)
		}
	} else {
		w.WriteInt(damage)
	}
	return w.Bytes()
}
