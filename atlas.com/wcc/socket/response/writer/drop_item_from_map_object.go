package writer

import (
   "atlas-wcc/socket/response"
   "time"
)

const OpCodeDropItemFromMapObject uint16 = 0x10C

func WriteDropItemFromMapObject(itemUniqueId uint32, itemId uint32, quantity uint32, meso uint32,
   dropperUniqueId uint32, dropType byte, ownerId uint32, ownerPartyId uint32, observerId uint32,
   observerPartyId uint32, dropTime uint64, dropX int16, dropY int16, dropperX int16, dropperY int16,
   characterDrop bool, mod byte) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeDropItemFromMapObject)

   ldt := dropType
   if hasClientSideOwnership(ownerId, ownerPartyId, observerId, observerPartyId, dropTime) && ldt < 3 {
      ldt = 2
   }

   w.WriteByte(mod)
   w.WriteInt(itemUniqueId)
   w.WriteBool(meso > 0) // 1 mesos, 0 item, 2 and above all item meso bag,
   w.WriteInt(itemId)
   w.WriteInt(getClientSideOwnerId(ownerId, ownerPartyId))
   // 0 = timeout for non-owner, 1 = timeout for non-owner's party, 2 = FFA, 3 = explosive/FFA
   w.WriteByte(dropType)
   w.WriteInt16(dropX)
   w.WriteInt16(dropY)
   w.WriteInt(dropperUniqueId)
   if mod != 2 {
      w.WriteInt16(dropperX)
      w.WriteInt16(dropperY)
      w.WriteShort(0)
   }
   if meso == 0 {
      //TODO add expiration (if necessary?)
      w.WriteInt64(-1)
   }
   if characterDrop {
      w.WriteByte(0)
   } else {
      w.WriteByte(1)
   }
   return w.Bytes()
}

func hasClientSideOwnership(ownerId uint32, ownerPartyId uint32, observerId uint32, observerPartyId uint32, dropTime uint64) bool {
   return ownerId == observerId ||
      (ownerPartyId != 0 && ownerPartyId == observerPartyId) ||
      hasExpiredOwnershipTime(dropTime)
}

func hasExpiredOwnershipTime(dropTime uint64) bool {
   return time.Now().UnixNano() - int64(dropTime) >= 15 * 1000
}

func getClientSideOwnerId(ownerId uint32, ownerPartyId uint32) uint32 {
   if ownerPartyId != 0 {
      return ownerPartyId
   }
   return ownerId
}
