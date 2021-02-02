package producers

import (
	"context"
	"log"
)

type characterAdjustHealthEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint16 `json:"amount"`
}

var CharacterAdjustHealth = func(l *log.Logger, ctx context.Context) *characterAdjustHealth {
	return &characterAdjustHealth{
		l:   l,
		ctx: ctx,
	}
}

type characterAdjustHealth struct {
	l   *log.Logger
	ctx context.Context
}

func (m *characterAdjustHealth) Emit(characterId uint32, amount uint16) {
	e := &characterAdjustHealthEvent{
		CharacterId: characterId,
		Amount: amount,
	}
	produceEvent(m.l, "TOPIC_ADJUST_HEALTH", createKey(int(characterId)), e)
}
