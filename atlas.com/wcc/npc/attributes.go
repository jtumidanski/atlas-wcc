package npc

import (
	"atlas-wcc/rest/response"
	"encoding/json"
)

type dataContainer struct {
	data response.DataSegment
}

type dataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
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

func (c *dataContainer) MarshalJSON() ([]byte, error) {
	t := struct {
		Data     interface{} `json:"data"`
		Included interface{} `json:"included"`
	}{}
	if len(c.data) == 1 {
		t.Data = c.data[0]
	} else {
		t.Data = c.data
	}
	return json.Marshal(t)
}

func (c *dataContainer) UnmarshalJSON(data []byte) error {
	d, _, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyNpcData))
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

func EmptyNpcData() interface{} {
	return &dataBody{}
}
