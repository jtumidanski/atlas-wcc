package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type characterDamageEvent struct {
	CharacterId     uint32 `json:"characterId"`
	MonsterId       uint32 `json:"monsterId"`
	MonsterUniqueId uint32 `json:"monsterUniqueId"`
	DamageFrom      int8   `json:"damageFrom"`
	Element         byte   `json:"element"`
	Damage          int32  `json:"damage"`
	Direction       int8   `json:"direction"`
}

var CharacterDamage = func(l logrus.FieldLogger, ctx context.Context) *characterDamage {
	return &characterDamage{
		l:   l,
		ctx: ctx,
	}
}

type characterDamage struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func (m *characterDamage) Emit(characterId uint32, monsterIdFrom uint32, uniqueId uint32, damageFrom int8, element byte, damage int32, direction int8) {
	e := &characterDamageEvent{
		CharacterId:     characterId,
		MonsterId:       monsterIdFrom,
		MonsterUniqueId: uniqueId,
		DamageFrom:      damageFrom,
		Element:         element,
		Damage:          damage,
		Direction:       direction,
	}
	produceEvent(m.l, "DAMAGE_CHARACTER", createKey(int(characterId)), e)
}
