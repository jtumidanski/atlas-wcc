package writer

import "atlas-wcc/socket/response"

const OpCodeCharacterHint uint16 = 0xD6

func WriteHint(message string, width int16, height int16) []byte {
   w := response.NewWriter()
   w.WriteShort(OpCodeCharacterHint)
   tw := width
   th := height
   if tw < 1 {
      tw = int16(len(message) * 10)
      if tw < 40 {
         tw = 40
      }
   }
   if th < 5 {
      th = 5
   }

   w.WriteAsciiString(message)
   w.WriteInt16(tw)
   w.WriteInt16(th)
   w.WriteByte(1)
   return w.Bytes()
}