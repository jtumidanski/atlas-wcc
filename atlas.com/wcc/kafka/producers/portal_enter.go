package producers

import (
	"context"
	"log"
)

type PortalEnter struct {
	l   *log.Logger
	ctx context.Context
}

func NewPortalEnter(l *log.Logger, ctx context.Context) *PortalEnter {
	return &PortalEnter{l, ctx}
}

func (m *PortalEnter) EmitEnter(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
	m.emit(worldId, channelId, mapId, portalId, characterId)
}

func (m *PortalEnter) emit(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
	e := &PortalEnterCommand{
		WorldId:     worldId,
		ChannelId:   channelId,
		MapId:       mapId,
		PortalId:    portalId,
		CharacterId: characterId,
	}
	ProduceEvent(m.l, "TOPIC_ENTER_PORTAL", createKey(int(characterId)), e)
}

type PortalEnterCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}
