package writer

import "atlas-wcc/socket/response"

const OpCodeRemoveNpc uint16 = 0x102

func WriteRemoveNPC(objectId uint32) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeRemoveNpc)
   w.WriteInt(objectId)
   return w.Bytes()
}