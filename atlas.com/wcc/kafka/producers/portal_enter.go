package producers

import (
	"context"
	"log"
)

type portalEnterCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

var PortalEnter = func(l *log.Logger, ctx context.Context) *portalEnter {
	return &portalEnter{
		l:   l,
		ctx: ctx,
	}
}

type portalEnter struct {
	l   *log.Logger
	ctx context.Context
}

func (m *portalEnter) Enter(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
	e := &portalEnterCommand{
		WorldId:     worldId,
		ChannelId:   channelId,
		MapId:       mapId,
		PortalId:    portalId,
		CharacterId: characterId,
	}
	produceEvent(m.l, "TOPIC_ENTER_PORTAL", createKey(int(characterId)), e)
}
