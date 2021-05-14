package shop

type Model struct {
	shopId uint32
	items  []Item
}

func (m Model) ShopId() uint32 {
	return m.shopId
}

func (m Model) Items() []Item {
	return m.items
}

type Item struct {
	itemId   uint32
	price    uint32
	pitch    uint32
	position uint32
}

func (i Item) ItemId() uint32 {
	return i.itemId
}

func (i Item) Price() uint32 {
	return i.price
}

func (i Item) Pitch() uint32 {
	return i.pitch
}
