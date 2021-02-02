package producers

import (
	"context"
	"log"
)

type CharacterStatus struct {
	l   *log.Logger
	ctx context.Context
}

func NewCharacterStatus(l *log.Logger, ctx context.Context) *CharacterStatus {
	return &CharacterStatus{l, ctx}
}

func (m *CharacterStatus) EmitLogin(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	m.emit(worldId, channelId, accountId, characterId, "LOGIN")
}

func (m *CharacterStatus) EmitLogout(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	m.emit(worldId, channelId, accountId, characterId, "LOGOUT")
}

func (m *CharacterStatus) emit(worldId byte, channelId byte, accountId uint32, characterId uint32, theType string) {
	e := &CharacterStatusEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		AccountId:   accountId,
		CharacterId: characterId,
		Type:        theType,
	}
	ProduceEvent(m.l, "TOPIC_CHARACTER_STATUS", createKey(int(characterId)), e)
}

type CharacterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}
