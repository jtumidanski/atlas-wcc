package attributes

import (
	"atlas-wcc/rest/response"
	"encoding/json"
)

type DamageDataContainer struct {
	data     response.DataSegment
	included response.DataSegment
}

type DamageData struct {
	Id         string           `json:"id"`
	Type       string           `json:"type"`
	Attributes DamageAttributes `json:"attributes"`
}

type DamageAttributes struct {
	Type    string `json:"type"`
	Maximum uint32 `json:"maximum"`
}

func (a *DamageDataContainer) MarshalJSON() ([]byte, error) {
	t := struct {
		Data     interface{} `json:"data"`
		Included interface{} `json:"included"`
	}{}
	if len(a.data) == 1 {
		t.Data = a.data[0]
	} else {
		t.Data = a.data
	}
	if len(a.included) == 1 {
		t.Included = a.included[0]
	} else {
		t.Included = a.included
	}
	return json.Marshal(t)
}

func (a *DamageDataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyDamageData))
	if err != nil {
		return err
	}

	a.data = d
	a.included = i
	return nil
}

func (a *DamageDataContainer) Data() *DamageData {
	if len(a.data) >= 1 {
		return a.data[0].(*DamageData)
	}
	return nil
}

func (a *DamageDataContainer) DataList() []DamageData {
	var r = make([]DamageData, 0)
	for _, x := range a.data {
		r = append(r, *x.(*DamageData))
	}
	return r
}

func EmptyDamageData() interface{} {
	return &DamageData{}
}
