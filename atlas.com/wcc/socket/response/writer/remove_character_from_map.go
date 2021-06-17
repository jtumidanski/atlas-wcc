package writer

import (
   "atlas-wcc/socket/response"
   "github.com/sirupsen/logrus"
)

const OpCodeRemoveCharacterFromMap uint16 = 0xA1

func WriteRemoveCharacterFromMap(l logrus.FieldLogger) func(characterId uint32) []byte {
   return func(characterId uint32) []byte {
   w := response.NewWriter(l)
   w.WriteShort(OpCodeRemoveCharacterFromMap)
   w.WriteInt(characterId)
   return w.Bytes()
   }
}
