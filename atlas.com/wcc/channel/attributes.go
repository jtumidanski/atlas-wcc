package channel

type attributes struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	Capacity  int    `json:"capacity"`
	IpAddress string `json:"ipAddress"`
	Port      uint16 `json:"port"`
}
