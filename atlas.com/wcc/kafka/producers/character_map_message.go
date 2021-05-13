package producers

import (
	"github.com/sirupsen/logrus"
)

type characterMapMessageEvent struct {
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	Message     string `json:"message"`
	GM          bool   `json:"gm"`
	Show        bool   `json:"show"`
}

func CharacterMapMessage(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, characterId uint32, message string, gm bool, show bool) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_MAP_MESSAGE_EVENT")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, message string, gm bool, show bool) {
		e := &characterMapMessageEvent{
			CharacterId: characterId,
			MapId:       mapId,
			Message:     message,
			GM:          gm,
			Show:        show,
		}
		producer(CreateKey(int(worldId)*1000+int(channelId)), e)
	}
}
