package domain

type Channel struct {
	worldId   byte
	channelId byte
	capacity  int
	ipAddress string
	port      uint16
}

func (c Channel) WorldId() byte {
	return c.worldId
}

func (c Channel) ChannelId() byte {
	return c.channelId
}

func (c Channel) Capacity() int {
	return c.capacity
}

func (c Channel) IpAddress() string {
	return c.ipAddress
}

func (c Channel) Port() uint16 {
	return c.port
}

type channelBuilder struct {
	worldId   byte
	channelId byte
	capacity  int
	ipAddress string
	port      uint16
}

func NewChannelBuilder() *channelBuilder {
	return &channelBuilder{}
}

func (c *channelBuilder) SetWorldId(worldId byte) *channelBuilder {
	c.worldId = worldId
	return c
}

func (c *channelBuilder) SetChannelId(channelId byte) *channelBuilder {
	c.channelId = channelId
	return c
}

func (c *channelBuilder) SetCapacity(capacity int) *channelBuilder {
	c.capacity = capacity
	return c
}

func (c *channelBuilder) SetIpAddress(ipAddress string) *channelBuilder {
	c.ipAddress = ipAddress
	return c
}

func (c *channelBuilder) SetPort(port uint16) *channelBuilder {
	c.port = port
	return c
}

func (c *channelBuilder) Build() Channel {
	return Channel{
		worldId:   c.worldId,
		channelId: c.channelId,
		capacity:  c.capacity,
		ipAddress: c.ipAddress,
		port:      c.port,
	}
}
