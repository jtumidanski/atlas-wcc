package producers

import (
	"context"
	"log"
)

type characterExpressionEvent struct {
	CharacterId uint32 `json:"characterId"`
	Emote       uint32 `json:"emote"`
}

var CharacterExpression = func(l *log.Logger, ctx context.Context) *characterExpression {
	return &characterExpression{
		l:   l,
		ctx: ctx,
	}
}

type characterExpression struct {
	l   *log.Logger
	ctx context.Context
}

func (m *characterExpression) Emit(characterId uint32, emote uint32) {
	e := &characterExpressionEvent{
		CharacterId: characterId,
		Emote:       emote,
	}
	produceEvent(m.l, "CHANGE_FACIAL_EXPRESSION", createKey(int(characterId)), e)
}
