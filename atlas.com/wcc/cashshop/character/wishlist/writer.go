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

		for _, wli := range wishlist {
			w.WriteInt(wli.SerialNumber())
		}
		for i := len(wishlist); i < 10; i++ {
			w.WriteInt(0)
		}
		return w.Bytes()
	}
}
