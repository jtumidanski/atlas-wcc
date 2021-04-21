package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type characterDistributeApEvent struct {
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

var CharacterDistributeAp = func(l logrus.FieldLogger, ctx context.Context) *characterDistributeAp {
	return &characterDistributeAp{
		l:   l,
		ctx: ctx,
	}
}

type characterDistributeAp struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func (m *characterDistributeAp) Emit(characterId uint32, attributeType string) {
	e := &characterDistributeApEvent{
		CharacterId: characterId,
		Type:        attributeType,
	}
	produceEvent(m.l, "TOPIC_ASSIGN_AP_COMMAND", createKey(int(characterId)), e)
}
