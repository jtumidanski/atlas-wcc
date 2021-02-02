package producers

import (
	"context"
	"log"
)

type characterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

var CharacterStatus = func(l *log.Logger, ctx context.Context) *characterStatus {
	return &characterStatus{
		l:   l,
		ctx: ctx,
	}
}

type characterStatus struct {
	l   *log.Logger
	ctx context.Context
}

func (m *characterStatus) Login(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	m.emit(worldId, channelId, accountId, characterId, "LOGIN")
}

func (m *characterStatus) Logout(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	m.emit(worldId, channelId, accountId, characterId, "LOGOUT")
}

func (m *characterStatus) emit(worldId byte, channelId byte, accountId uint32, characterId uint32, theType string) {
	e := &characterStatusEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		AccountId:   accountId,
		CharacterId: characterId,
		Type:        theType,
	}
	produceEvent(m.l, "TOPIC_CHARACTER_STATUS", createKey(int(characterId)), e)
}
