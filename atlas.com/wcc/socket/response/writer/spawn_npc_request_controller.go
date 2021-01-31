package writer

import (
   "atlas-wcc/domain"
   "atlas-wcc/socket/response"
)

const OpCodeSpawnNPCRequestController uint16 = 0x103

func WriteSpawnNPCController(npc domain.NPC, miniMap bool) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeSpawnNPCRequestController)
   w.WriteByte(1)
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
   w.WriteShort(npc.RX0())
   w.WriteShort(npc.RX1())
   w.WriteBool(miniMap)
   return w.Bytes()
}

func WriteRemoveNPCController(objectId uint32) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeSpawnNPCRequestController)
   w.WriteByte(0)
   w.WriteInt(objectId)
   return w.Bytes()
}
