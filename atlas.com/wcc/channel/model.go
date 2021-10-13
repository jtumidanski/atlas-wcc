package channel

type Model struct {
	worldId   byte
	channelId byte
	capacity  int
	ipAddress string
	port      uint16
}

func (c Model) WorldId() byte {
	return c.worldId
}

func (c Model) ChannelId() byte {
	return c.channelId
}

func (c Model) Capacity() int {
	return c.capacity
}

func (c Model) IpAddress() string {
	return c.ipAddress
}

func (c Model) Port() uint16 {
	return c.port
}

type builder struct {
	worldId   byte
	channelId byte
	capacity  int
	ipAddress string
	port      uint16
}

func NewBuilder() *builder {
	return &builder{}
}

func (c *builder) SetWorldId(worldId byte) *builder {
	c.worldId = worldId
	return c
}

func (c *builder) SetChannelId(channelId byte) *builder {
	c.channelId = channelId
	return c
}

func (c *builder) SetCapacity(capacity int) *builder {
	c.capacity = capacity
	return c
}

func (c *builder) SetIpAddress(ipAddress string) *builder {
	c.ipAddress = ipAddress
	return c
}

func (c *builder) SetPort(port uint16) *builder {
	c.port = port
	return c
}

func (c *builder) Build() Model {
	return Model{
		worldId:   c.worldId,
		channelId: c.channelId,
		capacity:  c.capacity,
		ipAddress: c.ipAddress,
		port:      c.port,
	}
}
