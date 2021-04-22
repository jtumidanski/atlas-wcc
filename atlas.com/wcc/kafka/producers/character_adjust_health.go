package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type characterAdjustHealthEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      int16  `json:"amount"`
}

var CharacterAdjustHealth = func(l logrus.FieldLogger, ctx context.Context) *characterAdjustHealth {
	return &characterAdjustHealth{
		l:   l,
		ctx: ctx,
	}
}

type characterAdjustHealth struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func (m *characterAdjustHealth) Emit(characterId uint32, amount int16) {
	e := &characterAdjustHealthEvent{
		CharacterId: characterId,
		Amount:      amount,
	}
	produceEvent(m.l, "TOPIC_ADJUST_HEALTH", createKey(int(characterId)), e)
}
