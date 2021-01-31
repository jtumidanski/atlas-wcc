package attributes

type AccountDataContainer struct {
	data     dataSegment
	included dataSegment
}

type AccountData struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes AccountAttributes `json:"attributes"`
}

type AccountAttributes struct {
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

func (a *AccountDataContainer) UnmarshalJSON(data []byte) error {
	d, i, err := unmarshalRoot(data, mapperFunc(EmptyAccountData))
	if err != nil {
		return err
	}

	a.data = d
	a.included = i
	return nil
}

func (a *AccountDataContainer) Data() *AccountData {
	if len(a.data) >= 1 {
		return a.data[0].(*AccountData)
	}
	return nil
}

func (a *AccountDataContainer) DataList() []AccountData {
	var r = make([]AccountData, 0)
	for _, x := range a.data {
		r = append(r, *x.(*AccountData))
	}
	return r
}

func EmptyAccountData() interface{} {
	return &AccountData{}
}
