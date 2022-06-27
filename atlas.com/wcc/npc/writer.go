package npc

import (
	"atlas-wcc/npc/shop"
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeSpawnNpc uint16 = 0x101
const OpCodeRemoveNpc uint16 = 0x102
const OpCodeSpawnNPCRequestController uint16 = 0x103
const OpCodeNPCAction uint16 = 0x104
const OpCodeNPCTalk uint16 = 0x130
const OpCodeOpenNPCShop uint16 = 0x131

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

func WriteNPCAnimation(l logrus.FieldLogger) func(objectId uint32, second byte, third byte) []byte {
	return func(objectId uint32, second byte, third byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCAction)
		w.WriteInt(objectId)
		w.WriteByte(second)
		w.WriteByte(third)
		return w.Bytes()
	}
}

func WriteNPCMove(l logrus.FieldLogger) func(movement []byte) []byte {
	return func(movement []byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCAction)
		w.WriteByteArray(movement)
		return w.Bytes()
	}
}

func WriteNPCTalk(l logrus.FieldLogger) func(npcId uint32, messageType byte, talk string, end []byte, speaker byte) []byte {
	return func(npcId uint32, messageType byte, talk string, end []byte, speaker byte) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCTalk)
		w.WriteByte(4)
		w.WriteInt(npcId)
		w.WriteByte(messageType)
		w.WriteByte(speaker)
		w.WriteAsciiString(talk)
		w.WriteByteArray(end)
		return w.Bytes()
	}
}

func WriteNPCTalkNum(l logrus.FieldLogger) func(npcId uint32, talk string, defaultValue int32, minimumValue int32, maximumValue int32) []byte {
	return func(npcId uint32, talk string, defaultValue int32, minimumValue int32, maximumValue int32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCTalk)
		w.WriteByte(4)
		w.WriteInt(npcId)
		w.WriteByte(3)
		w.WriteByte(0)
		w.WriteAsciiString(talk)
		w.WriteInt32(defaultValue)
		w.WriteInt32(minimumValue)
		w.WriteInt32(maximumValue)
		w.WriteInt(0)
		return w.Bytes()
	}
}

func WriteNPCTalkStyle(l logrus.FieldLogger) func(npcId uint32, talk string, styles []uint32) []byte {
	return func(npcId uint32, talk string, styles []uint32) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeNPCTalk)
		w.WriteByte(4)
		w.WriteInt(npcId)
		w.WriteByte(7)
		w.WriteByte(0)
		w.WriteAsciiString(talk)
		w.WriteByte(byte(len(styles)))
		for _, style := range styles {
			w.WriteInt(style)
		}
		return w.Bytes()
	}
}

func WriteGetNPCShop(l logrus.FieldLogger) func(model shop.Model) []byte {
	return func(model shop.Model) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeOpenNPCShop)
		w.WriteInt(model.ShopId())
		w.WriteShort(uint16(len(model.Items())))
		for _, i := range model.Items() {
			w.WriteInt(i.ItemId())
			w.WriteInt(i.Price())
			if i.Price() == 0 {
				w.WriteInt(i.Pitch())
			} else {
				w.WriteInt(0)
			}
			w.WriteInt(0) //Can be used x minutes after purchase
			w.WriteInt(0)

			//if (!ItemConstants.isRechargeable(item.itemId())) {
			w.WriteShort(1)
			w.WriteShort(1000)
			//	writer.writeShort(1);
			//	writer.writeShort(item.buyable());
			//} else {
			//	writer.writeShort(0);
			//	writer.writeInt(0);
			//	writer.writeShort(doubleToShortBits(ii.getUnitPrice(item.itemId())));
			//	writer.writeShort(ii.getSlotMax(packet.getClient(), item.itemId()));
			//}
		}
		return w.Bytes()
	}
}
