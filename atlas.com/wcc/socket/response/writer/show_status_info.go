package writer

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeShowStatusInfo uint16 = 0x27

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

func WriteShowInventoryFull(l logrus.FieldLogger) []byte {
	w := response.NewWriter(l)
	w.WriteShort(OpCodeShowStatusInfo)
	w.WriteByte(0)
	w.WriteByte(0xFF)
	w.WriteInt(0)
	w.WriteInt(0)
	return w.Bytes()
}
