package wishlist

type Model struct {
	serialNumber uint32
}

func (m Model) SerialNumber() uint32 {
	return m.serialNumber
}
