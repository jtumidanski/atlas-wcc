package writer

import "atlas-wcc/socket/response"

const OpCodeMoveCharacter uint16 = 0xB9

func WriteMoveCharacter(characterId uint32, movementList []byte) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeMoveCharacter)
	w.WriteInt(characterId)
	w.WriteInt(0)
	for _, b := range movementList {
		w.WriteByte(b)
	}
	return w.Bytes()
}
