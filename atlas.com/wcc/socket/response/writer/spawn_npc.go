package writer

import (
   "atlas-wcc/domain"
   "atlas-wcc/socket/response"
)

const OpCodeSpawnNpc uint16 = 0x101

func WriteSpawnNPC(npc domain.NPC) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeSpawnNpc)
   w.WriteInt(npc.ObjectId())
   w.WriteInt(npc.Id())
   w.WriteInt16(npc.X())
   w.WriteInt16(npc.CY())
   if npc.F() == 1 {
      w.WriteByte(0)
   } else {
      w.WriteByte(1)
   }
   w.WriteShort(npc.Fh())
   w.WriteInt16(npc.RX0())
   w.WriteInt16(npc.RX1())
   w.WriteByte(1)
   return w.Bytes()
}