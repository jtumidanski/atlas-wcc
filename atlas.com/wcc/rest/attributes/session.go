package attributes

type SessionListDataContainer struct {
	Data []SessionData `json:"data"`
}

type SessionData struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes SessionAttributes `json:"attributes"`
}

type SessionAttributes struct {
	AccountId uint32 `json:"accountId"`
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
}