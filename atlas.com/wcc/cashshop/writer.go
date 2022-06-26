package cashshop

import (
	"atlas-wcc/account"
	"atlas-wcc/character"
	"atlas-wcc/character/inventory"
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeSetCashShop uint16 = 0x7F
const OpCodeQueryCashResult uint16 = 0x144
const OpCodeCashShopOperation uint16 = 0x145

func WriteOpenCashShop(l logrus.FieldLogger) func(a account.Model, c character.Model, scis []SpecialCashItem) []byte {
	return func(a account.Model, c character.Model, scis []SpecialCashItem) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSetCashShop)
		character.AddCharacterInfo(w, c)
		w.WriteByte(1)
		w.WriteAsciiString(a.Name())
		w.WriteInt(0)
		w.WriteShort(uint16(len(scis)))
		for _, sci := range scis {
			w.WriteInt(sci.SN())
			w.WriteInt(sci.Modifier())
			w.WriteByte(sci.Info())
		}
		w.Skip(121)

		//TODO this needs to be reviewed to be dynamic. Not implemented in HeavenMS
		cd := []uint32{50200004, 50200069, 50200117, 50100008, 50000047}
		for i := uint32(1); i <= 8; i++ {
			for j := uint32(0); j < 2; j++ {
				for _, ci := range cd {
					w.WriteInt(i)
					w.WriteInt(j)
					w.WriteInt(ci)
				}
			}
		}

		w.WriteInt(0)
		w.WriteShort(0)
		w.WriteByte(0)
		w.WriteInt(75)
		return w.Bytes()
	}
}

func WriteCashInventory(l logrus.FieldLogger) func(a account.Model, c character.Model) []byte {
	return func(a account.Model, c character.Model) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCashShopOperation)
		w.WriteByte(0x4B)
		w.WriteShort(uint16(len(c.Inventory().CashInventory().Items())))
		for i, item := range c.Inventory().CashInventory().Items() {
			inventory.AddCashItemInformation(w, i, item, a.Id())
		}
		w.WriteShort(0) // character storage slots
		w.WriteInt16(a.CharacterSlots())
		return w.Bytes()
	}
}

func WriteCashGifts(l logrus.FieldLogger) func() []byte {
	return func() []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCashShopOperation)
		w.WriteByte(0x4D)
		w.WriteShort(0)
		// TODO load gifts
		return w.Bytes()
	}
}

func WriteCashInformation(l logrus.FieldLogger) func(credit uint32, points uint32, prepaid uint32) []byte {
	return func(credit uint32, points uint32, prepaid uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeQueryCashResult)
		w.WriteInt(credit)
		w.WriteInt(points)
		w.WriteInt(prepaid)
		return w.Bytes()
	}
}
