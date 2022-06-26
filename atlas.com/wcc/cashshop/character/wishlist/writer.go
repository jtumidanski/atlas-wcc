package wishlist

import (
	"atlas-wcc/socket/response"
	"github.com/sirupsen/logrus"
)

const OpCodeCashShopOperation uint16 = 0x145

func WriteWishList(l logrus.FieldLogger) func(wishlist []Model, update bool) []byte {
	return func(wishlist []Model, update bool) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeCashShopOperation)
		if update {
			w.WriteByte(0x55)
		} else {
			w.WriteByte(0x4F)
		}

		for i := 0; i < 10; i++ {
			w.WriteInt(wishlist[i].SerialNumber())
		}
		return w.Bytes()
	}
}
