package producers

import (
	"github.com/opentracing/opentracing-go"
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

func CharacterDamage(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, monsterIdFrom uint32, uniqueId uint32, damageFrom int8, element byte, damage int32, direction int8) {
	producer := ProduceEvent(l, span, "DAMAGE_CHARACTER")
	return func(characterId uint32, monsterIdFrom uint32, uniqueId uint32, damageFrom int8, element byte, damage int32, direction int8) {
		e := &characterDamageEvent{
			CharacterId:     characterId,
			MonsterId:       monsterIdFrom,
			MonsterUniqueId: uniqueId,
			DamageFrom:      damageFrom,
			Element:         element,
			Damage:          damage,
			Direction:       direction,
		}
		producer(CreateKey(int(characterId)), e)
	}
}
