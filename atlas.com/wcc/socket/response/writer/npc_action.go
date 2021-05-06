package writer

import "atlas-wcc/socket/response"

const OpCodeNPCAction uint16 = 0x104

func WriteNPCAnimation(objectId uint32, second byte, third byte) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeNPCAction)
	w.WriteInt(objectId)
	w.WriteByte(second)
	w.WriteByte(third)
	return w.Bytes()
}

func WriteNPCMove(movement []byte) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeNPCAction)
	w.WriteByteArray(movement)
	return w.Bytes()
}
