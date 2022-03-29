package properties

type attributes struct {
	AccountId          uint32 `json:"accountId"`
	WorldId            byte   `json:"worldId"`
	Name               string `json:"name"`
	Level              byte   `json:"level"`
	Experience         uint32 `json:"experience"`
	GachaponExperience uint32 `json:"gachaponExperience"`
	Strength           uint16 `json:"strength"`
	Dexterity          uint16 `json:"dexterity"`
	Intelligence       uint16 `json:"intelligence"`
	Luck               uint16 `json:"luck"`
	Hp                 uint16 `json:"hp"`
	MaxHp              uint16 `json:"maxHp"`
	Mp                 uint16 `json:"mp"`
	MaxMp              uint16 `json:"maxMp"`
	Meso               uint32 `json:"meso"`
	HpMpUsed           int    `json:"hpMpUsed"`
	JobId              uint16 `json:"jobId"`
	SkinColor          byte   `json:"skinColor"`
	Gender             byte   `json:"gender"`
	Fame               int16  `json:"fame"`
	Hair               uint32 `json:"hair"`
	Face               uint32 `json:"face"`
	Ap                 uint16 `json:"ap"`
	Sp                 string `json:"sp"`
	MapId              uint32 `json:"mapId"`
	SpawnPoint         byte   `json:"spawnPoint"`
	Gm                 int    `json:"gm"`
	X                  int16  `json:"x"`
	Y                  int16  `json:"y"`
	Stance             byte   `json:"stance"`
}
