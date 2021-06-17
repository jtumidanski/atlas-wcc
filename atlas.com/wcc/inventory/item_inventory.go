package inventory

type ItemInventory struct {
	capacity byte
	items    []Item
}

func NewItemInventory(capacity byte, items []Item) ItemInventory {
	return ItemInventory{
		capacity: capacity,
		items:    items,
	}
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
