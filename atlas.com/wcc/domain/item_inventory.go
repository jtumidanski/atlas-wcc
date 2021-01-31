package domain

type ItemInventory struct {
   capacity byte
   items []Item
}

func EmptyItemInventory() ItemInventory {
   return ItemInventory{
      capacity: 4,
      items:    make([]Item, 0),
   }
}

func (i ItemInventory) Capacity() byte {
   return i.capacity
}

func (i ItemInventory) Items() []Item {
   return i.items
}