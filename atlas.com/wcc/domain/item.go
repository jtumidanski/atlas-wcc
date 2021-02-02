package domain

type Item struct {
   itemId   uint32
   slot     int8
   quantity uint16
}

func NewItem(itemId uint32, slot int8, quantity uint16) Item {
   return Item{
      itemId:   itemId,
      slot:     slot,
      quantity: quantity,
   }
}

func (i Item) Slot() int8 {
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
