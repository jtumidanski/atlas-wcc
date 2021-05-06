package writer

import (
	"atlas-wcc/socket/response"
	"sort"
)

const OpCodeStatChanged uint16 = 0x1F

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

func WriteEnableActions() []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeStatChanged)
	w.WriteByte(1)
	w.WriteInt(0)
	return w.Bytes()
}

func WriteCharacterStatUpdate(updates []StatUpdate, enableActions bool) []byte {
	w := response.NewWriter()
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
