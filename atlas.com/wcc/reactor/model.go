package reactor

type Model struct {
	id             uint32
	classification uint32
	name           string
	state          int8
	eventState     byte
	delay          uint32
	direction      byte
	x              int16
	y              int16
	alive          bool
}

func (m Model) Id() uint32 {
	return m.id
}

func (m Model) Classification() uint32 {
	return m.classification
}

func (m Model) State() int8 {
	return m.state
}

func (m Model) X() int16 {
	return m.x
}

func (m Model) Y() int16 {
	return m.y
}

func (m Model) Alive() bool {
	return m.alive
}
