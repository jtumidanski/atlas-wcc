package shop

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
	NPC   uint32           `json:"npc"`
	Items []itemAttributes `json:"items"`
}

type itemAttributes struct {
	ItemId   uint32 `json:"itemId"`
	Price    uint32 `json:"price"`
	Pitch    uint32 `json:"pitch"`
	Position uint32 `json:"position"`
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
