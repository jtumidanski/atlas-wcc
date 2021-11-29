package member

type Model struct {
	id          uint32
	characterId uint32
	worldId     byte
	channelId   byte
	online      bool
}

func (m Model) WorldId() byte {
	return m.worldId
}

func (m Model) ChannelId() byte {
	return m.channelId
}

func (m Model) CharacterId() uint32 {
	return m.characterId
}
