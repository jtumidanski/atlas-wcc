package writer

import (
	"atlas-wcc/monster"
	"atlas-wcc/socket/response"
)

const OpCodeSpawnMonsterControl uint16 = 0xEE

func WriteControlMonster(m *monster.Model, newSpawn bool, aggro bool) []byte {
	w := response.NewWriter()
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

func WriteStopControlMonster(m *monster.Model) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeSpawnMonsterControl)
	w.WriteByte(0)
	w.WriteInt(m.UniqueId())
	return w.Bytes()
}
