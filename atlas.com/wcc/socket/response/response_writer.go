package response

import (
   "bytes"
   "encoding/binary"
   "log"
)

type Writer struct {
   o *bytes.Buffer
}

func NewWriter() *Writer {
   return &Writer{new(bytes.Buffer)}
}

// WriteInt16 -
func (w *Writer) WriteInt16(data int16) { w.WriteShort(uint16(data)) }

// WriteInt32 -
func (w *Writer) WriteInt32(data int32) { w.WriteInt(uint32(data)) }

// WriteInt64 -
func (w *Writer) WriteInt64(data int64) { w.WriteLong(uint64(data)) }

func (w *Writer) WriteInt(val uint32) {
   err := binary.Write(w.o, binary.LittleEndian, val)
   if err != nil {
      log.Fatal("[ERROR] writing int value")
   }
}

func (w *Writer) WriteShort(val uint16) {
   err := binary.Write(w.o, binary.LittleEndian, val)
   if err != nil {
      log.Fatal("[ERROR] writing short value")
   }
}

func (w *Writer) WriteLong(val uint64) {
   err := binary.Write(w.o, binary.LittleEndian, val)
   if err != nil {
      log.Fatal("[ERROR] writing long value")
   }
}

func (w *Writer) WriteByte(val byte) {
   err := binary.Write(w.o, binary.LittleEndian, val)
   if err != nil {
      log.Fatal("[ERROR] writing byte value")
   }
}

func (w *Writer) WriteByteArray(vals []byte) {
   for i := 0; i < len(vals); i++ {
      err := binary.Write(w.o, binary.LittleEndian, vals[i])
      if err != nil {
         log.Fatal("[ERROR] writing byte value")
      }
   }
}

func (w *Writer) WriteBool(val bool) {
   i := 1
   if !val {
      i = 0
   }
   w.WriteByte(byte(i))
}

func (w *Writer) WriteAsciiString(s string) {
   w.WriteShort(uint16(len(s)))
   w.WriteByteArray([]byte(s))
}

func (w *Writer) WriteKeyValue(key byte, value uint32) {
   w.WriteByte(key)
   w.WriteInt(value)
}

func (w *Writer) Bytes() []byte {
   return w.o.Bytes()
}
