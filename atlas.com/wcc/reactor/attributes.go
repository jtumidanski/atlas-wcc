package reactor

import "atlas-wcc/rest/response"

type dataContainer struct {
	data response.DataSegment
}

type dataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
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
	Alive           bool   `json:"alive"`
}

func (c *dataContainer) UnmarshalJSON(data []byte) error {
	d, _, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyData))
	if err != nil {
		return err
	}
	c.data = d
	return nil
}

func (c *dataContainer) Data() *dataBody {
	if len(c.data) >= 1 {
		return c.data[0].(*dataBody)
	}
	return nil
}

func (c *dataContainer) DataList() []dataBody {
	var r = make([]dataBody, 0)
	for _, x := range c.data {
		r = append(r, *x.(*dataBody))
	}
	return r
}

func EmptyData() interface{} {
	return &dataBody{}
}
