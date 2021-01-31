package domain

type Item struct {
   itemId   uint32
   slot     byte
   quantity uint16
}

func (i Item) Slot() byte {
   return i.slot
}

func (i Item) ItemId() uint32 {
   return i.itemId
}

func (i Item) Expiration() int64 {
   return 0
}

func (i Item) Quantity() uint16 {
   return i.quantity
}

func (i Item) Owner() string {
   return ""
}

func (i Item) Flag() uint16 {
   return 0
}
