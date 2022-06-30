package npc

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeSpawnNpc uint16 = 0x101
const OpCodeRemoveNpc uint16 = 0x102
const OpCodeSpawnNPCRequestController uint16 = 0x103

func WriteSpawnNPC(l logrus.FieldLogger) func(npc Model) []byte {
	return func(npc Model) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeSpawnNpc)
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
		w.WriteByte(1)
		return w.Bytes()
	}
}

func WriteRemoveNPC(l logrus.FieldLogger) func(objectId uint32) []byte {
	return func(objectId uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeRemoveNpc)
		w.WriteInt(objectId)
		return w.Bytes()
	}
}

func WriteSpawnNPCController(l logrus.FieldLogger) func(npc Model, miniMap bool) []byte {
	return func(npc Model, miniMap bool) []byte {
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
