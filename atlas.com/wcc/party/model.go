package party

type Model struct {
	id       uint32
	leaderId uint32
}

func (m Model) Id() uint32 {
	return m.id
}
