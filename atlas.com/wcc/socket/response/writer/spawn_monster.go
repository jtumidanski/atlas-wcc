package writer

import (
	"atlas-wcc/monster"
	"atlas-wcc/socket/response"
)

const OpCodeSpawnMonster uint16 = 0xEC

func WriteSpawnMonster(m monster.Model, newSpawn bool) []byte  {
   return WriteSpawnMonsterWithEffect(m, newSpawn, 0)
}

func WriteSpawnMonsterWithEffect(m monster.Model, newSpawn bool, effect byte) []byte {
   w := response.NewWriter()
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