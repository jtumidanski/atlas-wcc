package monster

type damageEntry struct {
	CharacterId uint32 `json:"characterId"`
	Damage      int64  `json:"damage"`
}

type attributes struct {
	WorldId            byte          `json:"worldId"`
	ChannelId          byte          `json:"channelId"`
	MapId              uint32        `json:"mapId"`
	MonsterId          uint32        `json:"monsterId"`
	ControlCharacterId uint32        `json:"controlCharacterId"`
	X                  int16         `json:"x"`
	Y                  int16         `json:"y"`
	FH                 int16         `json:"fh"`
	Stance             byte          `json:"stance"`
	Team               int8          `json:"team"`
	HP                 uint32        `json:"hp"`
	DamageEntries      []damageEntry `json:"damageEntries"`
}
