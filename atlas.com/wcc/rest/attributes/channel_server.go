package attributes

type ChannelServerDataContainer struct {
	data dataSegment
}

type ChannelServerData struct {
	Id         string                  `json:"id"`
	Type       string                  `json:"type"`
	Attributes ChannelServerAttributes `json:"attributes"`
}

type ChannelServerAttributes struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	Capacity  int    `json:"capacity"`
	IpAddress string `json:"ipAddress"`
	Port      uint16 `json:"port"`
}

func (c *ChannelServerDataContainer) UnmarshalJSON(data []byte) error {
	d, _, err := unmarshalRoot(data, mapperFunc(EmptyChannelServerData))
	if err != nil {
		return err
	}

	c.data = d
	return nil
}

func (c *ChannelServerDataContainer) Data() *ChannelServerData {
	if len(c.data) >= 1 {
		return c.data[0].(*ChannelServerData)
	}
	return nil
}

func (c *ChannelServerDataContainer) DataList() []ChannelServerData {
	var r = make([]ChannelServerData, 0)
	for _, x := range c.data {
		r = append(r, *x.(*ChannelServerData))
	}
	return r
}

func EmptyChannelServerData() interface{} {
	return &ChannelServerData{}
}
