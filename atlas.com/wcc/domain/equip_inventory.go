package domain

type EquipInventory struct {
   capacity byte
   items []EquippedItem
}

func EmptyEquipInventory() EquipInventory {
   return EquipInventory{
      capacity: 4,
      items:    make([]EquippedItem, 0),
   }
}

func (e EquipInventory) Capacity() byte {
   return e.capacity
}

func (e EquipInventory) Items() []EquippedItem {
   return e.items
}