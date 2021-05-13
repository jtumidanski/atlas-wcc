package producers

import "github.com/sirupsen/logrus"

type changeKeyMapCommand struct {
	CharacterId uint32         `json:"characterId"`
	Changes     []KeyMapChange `json:"changes"`
}

type KeyMapChange struct {
	Key        int32 `json:"key"`
	ChangeType int8  `json:"changeType"`
	Action     int32 `json:"action"`
}

func ChangeKeyMap(l logrus.FieldLogger) func(characterId uint32, changes []KeyMapChange) {
	producer := ProduceEvent(l, "TOPIC_CHANGE_KEY_MAP")
	return func(characterId uint32, changes []KeyMapChange) {
		e := changeKeyMapCommand{CharacterId: characterId, Changes: changes}
		producer(CreateKey(int(characterId)), e)
	}
}
