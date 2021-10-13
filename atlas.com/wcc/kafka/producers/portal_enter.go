package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type portalEnterCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

func PortalEnter(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_ENTER_PORTAL")
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