package skill

import "atlas-wcc/rest/response"

type DataContainer struct {
	data     response.DataSegment
	included response.DataSegment
}

type DataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  int64  `json:"expiration"`
}

func (a *DataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptySkillData))
	if err != nil {
		return err
	}

	a.data = d
	a.included = i
	return nil
}

func (a *DataContainer) Data() *DataBody {
	if len(a.data) >= 1 {
		return a.data[0].(*DataBody)
	}
	return nil
}

func (a *DataContainer) DataList() []DataBody {
	var r = make([]DataBody, 0)
	for _, x := range a.data {
		r = append(r, *x.(*DataBody))
	}
	return r
}

func EmptySkillData() interface{} {
	return &DataBody{}
}
