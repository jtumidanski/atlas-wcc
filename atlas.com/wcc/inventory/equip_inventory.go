package inventory

type EquipInventory struct {
	capacity byte
	items    []EquippedItem
}

func NewEquipInventory(capacity byte, items []EquippedItem) EquipInventory {
	return EquipInventory{capacity: capacity, items: items}
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
