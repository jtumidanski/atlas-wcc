package properties

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
	"sort"
)

const OpCodeStatChanged uint16 = 0x1F
const OpCodeShowStatusInfo uint16 = 0x27

const (
	StatUpdateSkin               uint32 = 0x1
	StatUpdateFace               uint32 = 0x2
	StatUpdateHair               uint32 = 0x4
	StatUpdateLevel              uint32 = 0x10
	StatUpdateJob                uint32 = 0x20
	StatUpdateStrength           uint32 = 0x40
	StatUpdateDexterity          uint32 = 0x80
	StatUpdateIntelligence       uint32 = 0x100
	StatUpdateLuck               uint32 = 0x200
	StatUpdateHP                 uint32 = 0x400
	StatUpdateMaxHP              uint32 = 0x800
	StatUpdateMP                 uint32 = 0x1000
	StatUpdateMaxMP              uint32 = 0x2000
	StatUpdateAvailableAP        uint32 = 0x4000
	StatUpdateAvailableSP        uint32 = 0x8000
	StatUpdateExperience         uint32 = 0x10000
	StatUpdateFame               uint32 = 0x20000
	StatUpdateMeso               uint32 = 0x40000
	StatUpdatePet                uint32 = 0x180008
	StatUpdateGachaponExperience uint32 = 0x200000
)

type StatUpdate struct {
	Stat   uint32
	Amount uint32
}

func NewStatUpdate(stat uint32, amount uint32) StatUpdate {
	return StatUpdate{
		Stat:   stat,
		Amount: amount,
	}
}

func WriteEnableActions(l logrus.FieldLogger) []byte {
	w := response.NewWriter(l)
	w.WriteShort(OpCodeStatChanged)
	w.WriteByte(1)
	w.WriteInt(0)
	return w.Bytes()
}

func WriteCharacterStatUpdate(l logrus.FieldLogger) func(updates []StatUpdate, enableActions bool) []byte {
	return func(updates []StatUpdate, enableActions bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeStatChanged)
		if enableActions {
			w.WriteByte(1)
		} else {
			w.WriteByte(0)
		}

		updateMask := uint32(0)
		for _, u := range updates {
			updateMask |= u.Stat
		}
		sortedUpdates := updates
		sort.SliceStable(sortedUpdates, func(i, j int) bool {
			return sortedUpdates[i].Stat < sortedUpdates[j].Stat
		})

		w.WriteInt(updateMask)
		for _, u := range sortedUpdates {
			if u.Stat >= 1 {
				if u.Stat == 0x1 {
					w.WriteByte(byte(u.Amount))
				} else if u.Stat <= 0x4 {
					w.WriteInt(u.Amount)
				} else if u.Stat < 0x20 {
					w.WriteByte(byte(u.Amount))
				} else if u.Stat == 0x8000 {
					w.WriteShort(uint16(u.Amount))
				} else if u.Stat < 0xFFFF {
					w.WriteShort(uint16(u.Amount))
				} else if u.Stat == 0x20000 {
					w.WriteShort(uint16(u.Amount))
				} else {
					w.WriteInt(u.Amount)
				}
			}
		}
		return w.Bytes()
	}
}

func WriteShowExperienceGain(l logrus.FieldLogger) func(gain uint32, equip uint32, party uint32, chat bool, white bool) []byte {
	return func(gain uint32, equip uint32, party uint32, chat bool, white bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowStatusInfo)
		w.WriteByte(3) // 3 = exp, 4 = fame, 5 = mesos, 6 = guild points
		w.WriteBool(white)
		w.WriteInt(gain)
		w.WriteBool(chat)
		w.WriteInt(0)  // bonus event exp
		w.WriteByte(0) // third monster kill event
		w.WriteByte(0)
		w.WriteInt(0) // wedding bonus
		if chat {
			w.WriteByte(0)
		}
		w.WriteByte(0) //0 = party bonus, 100 = 1x Bonus EXP, 200 = 2x Bonus EXP
		w.WriteInt(party)
		w.WriteInt(equip)
		w.WriteInt(0) // internet cafe
		w.WriteInt(0) // rainbow week
		return w.Bytes()
	}
}

func WriteShowMesoGain(l logrus.FieldLogger) func(gain int32, chat bool) []byte {
	return func(gain int32, chat bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowStatusInfo)
		if chat {
			w.WriteByte(5)
		} else {
			w.WriteByte(0)
			w.WriteShort(1)
		}
		w.WriteInt32(gain)
		w.WriteShort(0)
		return w.Bytes()
	}
}

func WriteShowItemGain(l logrus.FieldLogger) func(itemId uint32, quantity uint32) []byte {
	return func(itemId uint32, quantity uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeShowStatusInfo)
		w.WriteShort(0)
		w.WriteInt(itemId)
		w.WriteInt(quantity)
		w.WriteInt(0)
		w.WriteInt(0)
		return w.Bytes()
	}
}
