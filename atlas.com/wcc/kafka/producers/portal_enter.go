package producers

import (
	"github.com/sirupsen/logrus"
)

type portalEnterCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

func PortalEnter(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_ENTER_PORTAL")
	return func(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
		e := &portalEnterCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			MapId:       mapId,
			PortalId:    portalId,
			CharacterId: characterId,
		}
		producer(CreateKey(int(characterId)), e)
	}
}