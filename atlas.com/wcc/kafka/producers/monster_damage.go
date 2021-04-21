package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type monsterDamageEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	UniqueId    uint32 `json:"uniqueId"`
	CharacterId uint32 `json:"characterId"`
	Damage      uint32 `json:"damage"`
}

var MonsterDamage = func(l logrus.FieldLogger, ctx context.Context) *monsterDamage {
	return &monsterDamage{
		l:   l,
		ctx: ctx,
	}
}

type monsterDamage struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func (m *monsterDamage) Emit(worldId byte, channelId byte, mapId uint32, uniqueId uint32, characterId uint32, damage uint32) {
	e := &monsterDamageEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		MapId:       mapId,
		UniqueId:    uniqueId,
		CharacterId: characterId,
		Damage:      damage,
	}
	produceEvent(m.l, "TOPIC_MONSTER_DAMAGE", createKey(int(uniqueId)), e)
}
