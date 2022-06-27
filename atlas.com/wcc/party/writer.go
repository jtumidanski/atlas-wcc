package party

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodePartyOperation uint16 = 0x3E

func WritePartyCreated(l logrus.FieldLogger) func(id uint32) []byte {
	return func(id uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodePartyOperation)
		w.WriteInt(8)
		w.WriteInt(id)
		//TODO doors
		w.WriteInt(999999999)
		w.WriteInt(999999999)
		w.WriteInt(0)
		w.WriteInt(0)
		return w.Bytes()
	}
}

func WritePartyDisbanded(l logrus.FieldLogger) func(id uint32, characterId uint32) []byte {
	return func(id uint32, characterId uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodePartyOperation)
		w.WriteInt(0x0C)
		w.WriteInt(id)
		w.WriteInt(characterId)
		w.WriteInt(0)
		w.WriteInt(id)
		return w.Bytes()
	}
}
