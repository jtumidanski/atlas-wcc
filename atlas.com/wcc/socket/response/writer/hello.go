package writer

import (
	"atlas-wcc/socket/response"
)

func WriteHello(version uint16, sendIv []byte, recvIv []byte) []byte {
   w := response.NewWriter()
   w.WriteShort(uint16(0x0E))
   w.WriteShort(version)
   w.WriteShort(1)
   w.WriteByte(49)
   w.WriteByteArray(recvIv)
   w.WriteByteArray(sendIv)
   w.WriteByte(8)
   return w.Bytes()
}
