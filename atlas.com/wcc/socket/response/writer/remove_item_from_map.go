package writer

import "atlas-wcc/socket/response"

const OpCodeRemoveItemFromMap uint16 = 0x10D

func WriteRemoveItem(dropId uint32, animation byte, characterId uint32) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeRemoveItemFromMap)
   w.WriteByte(animation)
   w.WriteInt(dropId)
   if animation >= 2 {
      w.WriteInt(characterId)
      //TODO handle pet loot
   }
   return w.Bytes()
}
