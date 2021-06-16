package attributes

import "atlas-wcc/rest/response"

type NpcDataContainer struct {
	data response.DataSegment
}

type NpcData struct {
	Id         string        `json:"id"`
	Type       string        `json:"type"`
	Attributes NpcAttributes `json:"attributes"`
}

type NpcAttributes struct {
	Id   uint32 `json:"id"`
	Name string `json:"name"`
	CY   int16  `json:"cy"`
	F    uint32 `json:"f"`
	FH   uint16 `json:"fh"`
	RX0  int16  `json:"rx0"`
	RX1  int16  `json:"rx1"`
	X    int16  `json:"x"`
	Y    int16  `json:"y"`
	Hide bool   `json:"hide"`
}

func (c *NpcDataContainer) UnmarshalJSON(data []byte) error {
	d, _, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyNpcData))
	if err != nil {
		return err
	}
	c.data = d
	return nil
}

func (c *NpcDataContainer) Data() *NpcData {
	if len(c.data) >= 1 {
		return c.data[0].(*NpcData)
	}
	return nil
}

func (c *NpcDataContainer) DataList() []NpcData {
	var r = make([]NpcData, 0)
	for _, x := range c.data {
		r = append(r, *x.(*NpcData))
	}
	return r
}

func EmptyNpcData() interface{} {
	return &NpcData{}
}
