package producers

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type CharacterMovement struct {
	l   *log.Logger
	ctx context.Context
}

func NewCharacterMovement(l *log.Logger, ctx context.Context) *CharacterMovement {
	return &CharacterMovement{l, ctx}
}

func (m *CharacterMovement) EmitMovement(worldId byte, channelId byte, characterId uint32, x int16, y int16, stance byte, rawMovement []byte) {
	e := &CharacterMovementEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		CharacterId: characterId,
		X:           x,
		Y:           y,
		Stance:      stance,
		RawMovement: rawMovement,
	}
	ProduceEvent(m.l, "TOPIC_CHARACTER_MOVEMENT", createKey(int(characterId)), e)
}

type CharacterMovementEvent struct {
	WorldId     byte        `json:"worldId"`
	ChannelId   byte        `json:"channelId"`
	CharacterId uint32      `json:"characterId"`
	X           int16       `json:"x"`
	Y           int16       `json:"y"`
	Stance      byte        `json:"stance"`
	RawMovement RawMovement `json:"rawMovement"`
}

type RawMovement []byte

func (m RawMovement) MarshalJSON() ([]byte, error) {
	var result string
	if m == nil {
		result = "[]"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", m)), ",")
	}
	return []byte(result), nil
}
