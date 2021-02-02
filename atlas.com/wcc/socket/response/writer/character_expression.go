package writer

import "atlas-wcc/socket/response"

const OpCodeCharacterExpression uint16 = 0xC1

func WriteCharacterExpression(characterId uint32, expression uint32) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeCharacterExpression)
	w.WriteInt(characterId)
	w.WriteInt(expression)
	return w.Bytes()
}
