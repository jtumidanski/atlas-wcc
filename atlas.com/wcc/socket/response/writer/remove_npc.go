package writer

import (
   "atlas-wcc/socket/response"
   "github.com/sirupsen/logrus"
)

const OpCodeRemoveNpc uint16 = 0x102

func WriteRemoveNPC(l logrus.FieldLogger) func(objectId uint32) []byte {
   return func(objectId uint32) []byte {
   w := response.NewWriter(l)
   w.WriteShort(OpCodeRemoveNpc)
   w.WriteInt(objectId)
   return w.Bytes()
   }
}