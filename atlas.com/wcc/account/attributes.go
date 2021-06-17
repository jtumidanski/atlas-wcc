package account

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
	Name           string `json:"name"`
	Password       string `json:"password"`
	Pin            string `json:"pin"`
	Pic            string `json:"pic"`
	LoggedIn       int    `json:"loggedIn"`
	LastLogin      uint64 `json:"lastLogin"`
	Gender         byte   `json:"gender"`
	Banned         bool   `json:"banned"`
	TOS            bool   `json:"tos"`
	Language       string `json:"language"`
	Country        string `json:"country"`
	CharacterSlots int16  `json:"characterSlots"`
}

func (a *dataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := response.UnmarshalRoot(data, response.MapperFunc(EmptyAccountData))
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

func EmptyAccountData() interface{} {
	return &dataBody{}
}
