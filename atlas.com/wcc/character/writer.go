package character

import (
	"atlas-wcc/character/inventory"
	"atlas-wcc/pet"
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeCharacterHint uint16 = 0xD6

func WriteHint(l logrus.FieldLogger) func(message string, width int16, height int16) []byte {
	return func(message string, width int16, height int16) []byte {
		w := response.NewWriter(l)
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
}

func AddCharacterLook(w *response.Writer, character Model, mega bool) {
	w.WriteByte(character.Attributes().Gender())
	w.WriteByte(character.Attributes().SkinColor())
	w.WriteInt(character.Attributes().Face())
	if mega {
		w.WriteByte(0)
	} else {
		w.WriteByte(1)
	}
	w.WriteInt(character.Attributes().Hair())
	AddCharacterEquips(w, character)
}

func AddCharacterEquips(w *response.Writer, character Model) {
	var equips = getEquippedItemSlotMap(character.Equipment())
	var maskedEquips = make(map[int16]uint32)
	writeEquips(w, equips, maskedEquips)

	var weapon *inventory.EquippedItem
	for _, x := range character.Equipment() {
		if x.InWeaponSlot() {
			weapon = &x
			break
		}
	}
	if weapon != nil {
		w.WriteInt(weapon.ItemId())
	} else {
		w.WriteInt(0)
	}

	pet.WriteForEachPet(w, character.Pets(), pet.WritePetItemId, pet.WriteEmptyPetItemId)
}

func getEquippedItemSlotMap(e []inventory.EquippedItem) map[int16]uint32 {
	var equips = make(map[int16]uint32, 0)
	for _, x := range e {
		if x.NotInWeaponSlot() {
			y := x.InvertSlot()
			equips[y.Slot()] = y.ItemId()
		}
	}
	return equips
}

func writeEquips(w *response.Writer, equips map[int16]uint32, maskedEquips map[int16]uint32) {
	for k, v := range equips {
		w.WriteKeyValue(byte(k), v)
	}
	w.WriteByte(0xFF)
	for k, v := range maskedEquips {
		w.WriteKeyValue(byte(k), v)
	}
	w.WriteByte(0xFF)
}

func AddRingLook(w *response.Writer, _ Model, _ bool) {
	w.WriteByte(0)
}

func AddMarriageRingLook(w *response.Writer, _ Model, _ Model) {
	w.WriteByte(0)
}
