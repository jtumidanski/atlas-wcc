package producers

import (
	"context"
	"log"
)

type characterDistributeSpEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
}

var CharacterDistributeSp = func(l *log.Logger, ctx context.Context) *characterDistributeSp {
	return &characterDistributeSp{
		l:   l,
		ctx: ctx,
	}
}

type characterDistributeSp struct {
	l   *log.Logger
	ctx context.Context
}

func (m *characterDistributeSp) Emit(characterId uint32, skillId uint32) {
	e := &characterDistributeSpEvent{
		CharacterId: characterId,
		SkillId:     skillId,
	}
	produceEvent(m.l, "TOPIC_ASSIGN_SP_COMMAND", createKey(int(characterId)), e)
}
