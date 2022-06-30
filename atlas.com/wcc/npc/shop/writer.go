package shop

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeOpenNPCShop uint16 = 0x131

func WriteGetNPCShop(l logrus.FieldLogger) func(model Model) []byte {
	return func(model Model) []byte {
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
