package writer

import "atlas-wcc/socket/response"

const OpCodeRemoveCharacterFromMap uint16 = 0xA1

func WriteRemoveCharacterFromMap(characterId uint32) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeRemoveCharacterFromMap)
   w.WriteInt(characterId)
   return w.Bytes()
}
