package character

import (
	"atlas-wcc/character/inventory"
	"atlas-wcc/pet"
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
	"time"
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

func AddCharacterInfo(w *response.Writer, character Model) {
	w.WriteInt64(-1)
	w.WriteByte(0)
	addCharacterStats(w, character)
	//buddy list capacity
	w.WriteByte(0)
	//      if (character.getLinkedName() == null) {
	w.WriteByte(0)
	//      } else {
	//         writer.write(1);
	//         writer.writeMapleAsciiString(character.getLinkedName());
	//      }
	w.WriteInt(character.Attributes().Meso())
	inventory.AddInventoryInfo(w, character.Equipment(), character.Inventory())
	addSkillInfo(w, character)
	addQuestInfo(w, character)
	addMiniGameInfo(w, character)
	addRingInfo(w, character)
	addTeleportInfo(w, character)
	addMonsterBookInfo(w, character)
	addNewYearInfo(w, character)
	addAreaInfo(w, character)
	w.WriteShort(0)
}

func addAreaInfo(w *response.Writer, _ Model) {
	w.WriteShort(0)
}

func addNewYearInfo(w *response.Writer, _ Model) {
	w.WriteShort(0)
}

func addMonsterBookInfo(w *response.Writer, _ Model) {
	w.WriteInt(0)
	w.WriteByte(0)
	w.WriteShort(0)
}

func addTeleportInfo(w *response.Writer, _ Model) {
	for i := 0; i < 5; i++ {
		w.WriteInt(999999999)
	}
	for j := 0; j < 10; j++ {
		w.WriteInt(999999999)
	}
}

func addRingInfo(w *response.Writer, _ Model) {
	w.WriteShort(0)
	w.WriteShort(0)
	w.WriteShort(0)
}

func addMiniGameInfo(w *response.Writer, _ Model) {
	w.WriteShort(0)
}

func addQuestInfo(w *response.Writer, _ Model) {
	w.WriteShort(0)
	w.WriteShort(0)
}

func addSkillInfo(w *response.Writer, character Model) {
	w.WriteByte(0)

	sc := uint16(0)
	for _, s := range character.Skills() {
		if !s.Hidden() {
			sc += 1
		}
	}
	w.WriteShort(sc)

	for _, s := range character.Skills() {
		if !s.Hidden() {
			w.WriteInt(s.Id())
			w.WriteInt(s.Level())
			w.WriteLong(uint64(getTime(s.Expiration())))
			if s.FourthJob() {
				w.WriteInt(s.MasterLevel())
			}
		}
	}

	//      writer.writeShort(character.getAllCoolDowns().size());
	w.WriteShort(0)
	//      for (PlayerCoolDownValueHolder cooling : character.getAllCoolDowns()) {
	//         writer.writeInt(cooling.skillId);
	//         int timeLeft = (int) (cooling.length + cooling.startTime - System.currentTimeMillis());
	//         writer.writeShort(timeLeft / 1000);
	//      }
}

func addCharacterStats(w *response.Writer, character Model) {
	w.WriteInt(character.Attributes().Id())
	addPaddedCharacterName(w, character.Attributes().Name())
	w.WriteByte(character.Attributes().Gender())
	w.WriteByte(character.Attributes().SkinColor())
	w.WriteInt(character.Attributes().Face())
	w.WriteInt(character.Attributes().Hair())
	pet.WriteForEachPet(w, character.Pets(), writePetId, writeEmptyPetId)
	w.WriteByte(character.Attributes().Level())
	w.WriteShort(character.Attributes().JobId())
	w.WriteShort(character.Attributes().Strength())
	w.WriteShort(character.Attributes().Dexterity())
	w.WriteShort(character.Attributes().Intelligence())
	w.WriteShort(character.Attributes().Luck())
	w.WriteShort(character.Attributes().Hp())
	w.WriteShort(character.Attributes().MaxHp())
	w.WriteShort(character.Attributes().Mp())
	w.WriteShort(character.Attributes().MaxMp())
	w.WriteShort(character.Attributes().Ap())

	if character.Attributes().HasSPTable() {
		addRemainingSkillInfo(w, character)
	} else {
		w.WriteShort(character.Attributes().RemainingSp())
	}

	w.WriteInt(character.Attributes().Experience())
	w.WriteShort(uint16(character.Attributes().Fame()))
	w.WriteInt(character.Attributes().GachaponExperience())
	w.WriteInt(character.Attributes().MapId())
	w.WriteByte(character.Attributes().SpawnPoint())
	w.WriteInt(0)
}
func addRemainingSkillInfo(_ *response.Writer, _ Model) {

}

func addPaddedCharacterName(w *response.Writer, name string) {
	if len(name) > 13 {
		name = name[:13]
	}
	padSize := 13 - len(name)
	w.WriteByteArray([]byte(name))
	for i := 0; i < padSize; i++ {
		w.WriteByte(0x0)
	}
}

func writePetId(w *response.Writer, pet pet.Model) {
	w.WriteLong(pet.Id())
}

func writeEmptyPetId(w *response.Writer) {
	w.WriteLong(0)
}

func getTime(utcTimestamp int64) int64 {
	if utcTimestamp < 0 && utcTimestamp >= -3 {
		if utcTimestamp == -1 {
			return DefaultTime //high number ll
		} else if utcTimestamp == -2 {
			return ZeroTime
		} else {
			return Permanent
		}
	}

	ftUtOffset := 116444736010800000 + (10000 * timeNow())
	return utcTimestamp*10000 + ftUtOffset
}

func timeNow() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

const (
	DefaultTime int64 = 150842304000000000
	ZeroTime    int64 = 94354848000000000
	Permanent   int64 = 150841440000000000
)
