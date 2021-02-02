package producers

import (
	"context"
	"log"
)

type characterAdjustManaEvent struct {
	CharacterId uint32 `json:"characterId"`
	Amount      uint16 `json:"amount"`
}

var CharacterAdjustMana = func(l *log.Logger, ctx context.Context) *characterAdjustMana {
	return &characterAdjustMana{
		l:   l,
		ctx: ctx,
	}
}

type characterAdjustMana struct {
	l   *log.Logger
	ctx context.Context
}

func (m *characterAdjustMana) Emit(characterId uint32, amount uint16) {
	e := &characterAdjustManaEvent{
		CharacterId: characterId,
		Amount: amount,
	}
	produceEvent(m.l, "TOPIC_ADJUST_MANA", createKey(int(characterId)), e)
}
