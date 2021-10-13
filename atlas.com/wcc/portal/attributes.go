package portal

import "atlas-wcc/rest/response"

type dataContainer struct {
	data     response.DataSegment
	included response.DataSegment
}

type dataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
	Name        string `json:"name"`
	Target      string `json:"target"`
	Type        uint8  `json:"type"`
	X           int16  `json:"x"`
	Y           int16  `json:"y"`
	TargetMapId uint32 `json:"target_map_id"`
	ScriptName  string `json:"script_name"`
}

func (a *dataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyData))
	if err != nil {
		return err
	}

	a.data = d
	a.included = i
	return nil
}

func (a *dataContainer) Data() *dataBody {
	if len(a.data) >= 1 {
		return a.data[0].(*dataBody)
	}
	return nil
}

func (a *dataContainer) DataList() []dataBody {
	var r = make([]dataBody, 0)
	for _, x := range a.data {
		r = append(r, *x.(*dataBody))
	}
	return r
}

func EmptyData() interface{} {
	return &dataBody{}
}
