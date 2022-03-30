package monster

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeSpawnMonster uint16 = 0xEC
const OpCodeKillMonster uint16 = 0xED
const OpCodeSpawnMonsterControl uint16 = 0xEE
const OpCodeMoveMonster uint16 = 0xEF
const OpCodeMoveMonsterResponse uint16 = 0xF0

func WriteSpawnMonster(l logrus.FieldLogger) func(m *Model, newSpawn bool) []byte {
	return func(m *Model, newSpawn bool) []byte {
		return WriteSpawnMonsterWithEffect(l)(m, newSpawn, 0)
	}
}

func WriteSpawnMonsterWithEffect(l logrus.FieldLogger) func(m *Model, newSpawn bool, effect byte) []byte {
	return func(m *Model, newSpawn bool, effect byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSpawnMonster)
		w.WriteInt(m.UniqueId())
		if m.Controlled() {
			w.WriteByte(1)
		} else {
			w.WriteByte(5)
		}
		w.WriteInt(m.MonsterId())
		w.Skip(15)
		w.WriteByte(0x88)
		w.Skip(6)
		w.WriteInt16(m.X())
		w.WriteInt16(m.Y())
		w.WriteByte(m.Stance())
		w.WriteInt16(0)
		w.WriteInt16(m.FH())
		/*
		 * -4: Fake -3: Appear after linked mob is dead -2: Fade in 1: Smoke 3:
		 * King Slime spawn 4: Summoning rock thing, used for 3rd job? 6:
		 * Magical shit 7: Smoke shit 8: 'The Boss' 9/10: Grim phantom shit?
		 * 11/12: Nothing? 13: Frankenstein 14: Angry ^ 15: Orb animation thing,
		 * ?? 16: ?? 19: Mushroom castle boss thing
		 */

		//TODO handle parent mobs
		//      if (life.getParentMobOid() != 0) {
		//         MapleMonster parentMob = life.getMap().getMonsterByOid(life.getParentMobOid());
		//         if (parentMob != null && parentMob.isAlive()) {
		//            writer.write(packet.getEffect() != 0 ? packet.getEffect() : -3);
		//            writer.writeInt(life.getParentMobOid());
		//         } else {
		//            encodeParentlessMobSpawnEffect(writer, packet.isNewSpawn(), packet.getEffect());
		//         }
		//      } else {
		encodeParentlessMobSpawnEffect(w, newSpawn, effect)
		//      }

		w.WriteInt8(m.Team())
		w.WriteInt(0) // getItemEffect
		return w.Bytes()
	}
}

func encodeParentlessMobSpawnEffect(w *response.Writer, newSpawn bool, effect byte) {
	if effect > 0 {
		w.WriteByte(effect)
		w.WriteByte(0)
		w.WriteShort(0)
		if effect == 15 {
			w.WriteByte(0)
		}
	}
	if newSpawn {
		w.WriteInt8(-2)
	} else {
		w.WriteInt8(-1)
	}
}

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

func WriteControlMonster(l logrus.FieldLogger) func(m *Model, newSpawn bool, aggro bool) []byte {
	return func(m *Model, newSpawn bool, aggro bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSpawnMonsterControl)
		if aggro {
			w.WriteByte(2)
		} else {
			w.WriteByte(1)
		}
		w.WriteInt(m.UniqueId())
		if m.Controlled() {
			w.WriteByte(1)
		} else {
			w.WriteByte(5)
		}
		w.WriteInt(m.MonsterId())
		w.Skip(15)
		w.WriteByte(0x88)
		w.Skip(6)
		w.WriteInt16(m.X())
		w.WriteInt16(m.Y())
		w.WriteByte(m.Stance())
		w.WriteInt16(0)
		w.WriteInt16(m.FH())
		/*
		 * -4: Fake -3: Appear after linked mob is dead -2: Fade in 1: Smoke 3:
		 * King Slime spawn 4: Summoning rock thing, used for 3rd job? 6:
		 * Magical shit 7: Smoke shit 8: 'The Boss' 9/10: Grim phantom shit?
		 * 11/12: Nothing? 13: Frankenstein 14: Angry ^ 15: Orb animation thing,
		 * ?? 16: ?? 19: Mushroom castle boss thing
		 */

		//TODO handle parent mobs
		//      if (life.getParentMobOid() != 0) {
		//         MapleMonster parentMob = life.getMap().getMonsterByOid(life.getParentMobOid());
		//         if (parentMob != null && parentMob.isAlive()) {
		//            writer.write(packet.getEffect() != 0 ? packet.getEffect() : -3);
		//            writer.writeInt(life.getParentMobOid());
		//         } else {
		//            encodeParentlessMobSpawnEffect(writer, packet.isNewSpawn(), packet.getEffect());
		//         }
		//      } else {
		encodeParentlessMobSpawnEffect(w, newSpawn, 0)
		//      }

		w.WriteInt8(m.Team())
		w.WriteInt(0) // getItemEffect
		return w.Bytes()
	}
}

func WriteStopControlMonster(l logrus.FieldLogger) func(m *Model) []byte {
	return func(m *Model) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSpawnMonsterControl)
		w.WriteByte(0)
		w.WriteInt(m.UniqueId())
		return w.Bytes()
	}
}

func WriteMoveMonster(l logrus.FieldLogger) func(objectId uint32, skillPossible bool, skill int8, skillId uint32, skillLevel uint32, option uint16, startX int16, startY int16, movementList []byte) []byte {
	return func(objectId uint32, skillPossible bool, skill int8, skillId uint32, skillLevel uint32, option uint16, startX int16, startY int16, movementList []byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeMoveMonster)
		w.WriteInt(objectId)
		w.WriteByte(0)
		w.WriteBool(skillPossible)
		w.WriteInt8(skill)
		w.WriteByte(byte(skillId))
		w.WriteByte(byte(skillLevel))
		w.WriteShort(option)
		w.WriteInt16(startX)
		w.WriteInt16(startY)
		for _, b := range movementList {
			w.WriteByte(b)
		}
		return w.Bytes()
	}
}

func WriteMoveMonsterResponse(l logrus.FieldLogger) func(objectId uint32, moveId uint16, currentMp uint16, useSkills bool, skillId byte, skillLevel byte) []byte {
	return func(objectId uint32, moveId uint16, currentMp uint16, useSkills bool, skillId byte, skillLevel byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeMoveMonsterResponse)
		w.WriteInt(objectId)
		w.WriteShort(moveId)
		w.WriteBool(useSkills)
		w.WriteShort(currentMp)
		w.WriteByte(skillId)
		w.WriteByte(skillLevel)
		return w.Bytes()
	}
}
