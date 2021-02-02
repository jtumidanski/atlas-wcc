package attributes

type DamageDataContainer struct {
	data     dataSegment
	included dataSegment
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

func (a *DamageDataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := unmarshalRoot(data, mapperFunc(EmptyDamageData))
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
