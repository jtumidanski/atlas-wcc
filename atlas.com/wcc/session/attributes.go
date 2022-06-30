package session

type DataListContainer struct {
	Data []DataBody `json:"data"`
}

type DataBody struct {
	Id         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes Attributes `json:"attributes"`
}

type Attributes struct {
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
}
