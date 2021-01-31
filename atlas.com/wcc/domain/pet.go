package domain

type Pet struct {
	id     uint64
	itemId uint32
}

func (p Pet) Id() uint64 {
	return p.id
}

func (p Pet) ItemId() uint32 {
	return p.itemId
}
