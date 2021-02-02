package writer

import "atlas-wcc/socket/response"

const OpCodeCloseRangeAttack uint16 = 0xBA

func WriteCloseRangeAttack(characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeCloseRangeAttack)
	addAttackBody(w, characterId, skill, skillLevel, stance, numberAttackedAndDamaged, 0, damage, speed, direction, display)
	return w.Bytes()
}

func addAttackBody(w *response.Writer, characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, projectile uint32, damage map[uint32][]uint32, speed byte, direction byte, display byte) {
	w.WriteInt(characterId)
	w.WriteByte(numberAttackedAndDamaged)
	w.WriteByte(0x5b)
	w.WriteByte(skillLevel)
	if skillLevel > 0 {
		w.WriteInt(skill)
	}
	w.WriteByte(display)
	w.WriteByte(direction)
	w.WriteByte(stance)
	w.WriteByte(speed)
	w.WriteByte(0x0A)
	w.WriteInt(projectile)
	for k, v := range damage {
		w.WriteInt(k)
		w.WriteByte(0x0)
		if skill == 4211006 {
			w.WriteByte(byte(len(v)))
		}
		for _, e := range v {
			w.WriteInt(e)
		}
	}
}
