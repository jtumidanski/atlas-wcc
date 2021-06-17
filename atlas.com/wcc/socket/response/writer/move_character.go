package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeMoveCharacter uint16 = 0xB9

func WriteMoveCharacter(l logrus.FieldLogger) func(characterId uint32, movementList []byte) []byte {
	return func(characterId uint32, movementList []byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeMoveCharacter)
		w.WriteInt(characterId)
		w.WriteInt(0)
		for _, b := range movementList {
			w.WriteByte(b)
		}
		return w.Bytes()
	}
}
