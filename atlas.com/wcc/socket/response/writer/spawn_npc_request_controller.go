package writer

import (
	"atlas-wcc/npc"
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeSpawnNPCRequestController uint16 = 0x103

func WriteSpawnNPCController(l logrus.FieldLogger) func(npc npc.Model, miniMap bool) []byte {
	return func(npc npc.Model, miniMap bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSpawnNPCRequestController)
		w.WriteByte(1)
		w.WriteInt(npc.ObjectId())
		w.WriteInt(npc.Id())
		w.WriteInt16(npc.X())
		w.WriteInt16(npc.CY())
		if npc.F() == 1 {
			w.WriteByte(0)
		} else {
			w.WriteByte(1)
		}
		w.WriteShort(npc.Fh())
		w.WriteInt16(npc.RX0())
		w.WriteInt16(npc.RX1())
		w.WriteBool(miniMap)
		return w.Bytes()
	}
}

func WriteRemoveNPCController(l logrus.FieldLogger) func(objectId uint32) []byte {
	return func(objectId uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSpawnNPCRequestController)
		w.WriteByte(0)
		w.WriteInt(objectId)
		return w.Bytes()
	}
}
