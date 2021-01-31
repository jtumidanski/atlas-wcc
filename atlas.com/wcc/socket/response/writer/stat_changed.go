package writer

import (
   "atlas-wcc/socket/response"
)

const OpCodeStatChanged uint16 = 0x1F

func WriteEnableActions() []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeStatChanged)
   w.WriteByte(1)
   w.WriteInt(0)
   return w.Bytes()
}