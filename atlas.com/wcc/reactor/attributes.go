package reactor

type DataContainer struct {
	Data DataBody `json:"data"`
}

type DataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	WorldId         byte   `json:"world_id"`
	ChannelId       byte   `json:"channel_id"`
	MapId           uint32 `json:"map_id"`
	Classification  uint32 `json:"classification"`
	Name            string `json:"name"`
	Type            int32  `json:"type"`
	State           int8   `json:"state"`
	EventState      byte   `json:"event_state"`
	X               int16  `json:"x"`
	Y               int16  `json:"y"`
	Delay           uint32 `json:"delay"`
	FacingDirection byte   `json:"facing_direction"`
}