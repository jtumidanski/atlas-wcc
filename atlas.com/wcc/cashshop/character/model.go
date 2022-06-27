package character

type Model struct {
	characterId uint32
	credit      uint32
	points      uint32
	prepaid     uint32
}

func (m *Model) Credit() uint32 {
	return m.credit
}

func (m *Model) Points() uint32 {
	return m.points
}

func (m *Model) Prepaid() uint32 {
	return m.prepaid
}

func (m *Model) CharacterId() uint32 {
	return m.characterId
}
