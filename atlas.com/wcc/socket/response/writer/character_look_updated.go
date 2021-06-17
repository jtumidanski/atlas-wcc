package writer

import (
	"atlas-wcc/character"
	"atlas-wcc/socket/response"
)

const OpCodeUpdateCharacterLook uint16 = 0xC5

func WriteCharacterLookUpdated(r character.Model, c character.Model) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeUpdateCharacterLook)
	w.WriteInt(c.Attributes().Id())
	w.WriteByte(1)
	addCharacterLook(w, c, false)
	addRingLook(w, c, true)
	addRingLook(w, c, false)
	addMarriageRingLook(w, r, c)
	w.WriteInt(0)
	return w.Bytes()
}
