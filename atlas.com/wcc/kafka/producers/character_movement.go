package producers

import (
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strings"
)

type characterMovementEvent struct {
	WorldId     byte        `json:"worldId"`
	ChannelId   byte        `json:"channelId"`
	CharacterId uint32      `json:"characterId"`
	X           int16       `json:"x"`
	Y           int16       `json:"y"`
	Stance      byte        `json:"stance"`
	RawMovement rawMovement `json:"rawMovement"`
}

func MoveCharacter(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32, x int16, y int16, stance byte, rawMovement []byte) {
	producer := ProduceEvent(l, span, "TOPIC_CHARACTER_MOVEMENT")
	return func(worldId byte, channelId byte, characterId uint32, x int16, y int16, stance byte, rawMovement []byte) {
		e := &characterMovementEvent{
			WorldId:     worldId,
			ChannelId:   channelId,
			CharacterId: characterId,
			X:           x,
			Y:           y,
			Stance:      stance,
			RawMovement: rawMovement,
		}
		producer(CreateKey(int(characterId)), e)
	}
}

type rawMovement []byte

func (m rawMovement) MarshalJSON() ([]byte, error) {
	var result string
	if m == nil {
		result = "[]"
	} else {
		result = strings.Join(strings.Fields(fmt.Sprintf("%d", m)), ",")
	}
	return []byte(result), nil
}
