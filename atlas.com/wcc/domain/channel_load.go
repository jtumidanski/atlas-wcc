package domain

type ChannelLoad struct {
	channelId byte
	capacity  int
}

func NewChannelLoad(channelId byte, capacity int) ChannelLoad {
	return ChannelLoad{channelId, capacity}
}

func (cl ChannelLoad) ChannelId() byte {
	return cl.channelId
}

func (cl ChannelLoad) Capacity() int {
	return cl.capacity
}