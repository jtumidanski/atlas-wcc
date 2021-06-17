package inventory

type Item struct {
   itemId   uint32
   slot     int16
   quantity uint16
}

func NewItem(itemId uint32, slot int16, quantity uint16) Item {
   return Item{
      itemId:   itemId,
      slot:     slot,
      quantity: quantity,
   }
}

func (i Item) Slot() int16 {
   return i.slot
}

func (i Item) ItemId() uint32 {
   return i.itemId
}

func (i Item) Expiration() int64 {
   return -1
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
