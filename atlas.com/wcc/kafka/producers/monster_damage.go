package producers

import (
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

func MonsterDamage(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, characterId uint32, damage uint32) {
	producer := ProduceEvent(l, "TOPIC_MONSTER_DAMAGE")
	return func(worldId byte, channelId byte, mapId uint32, uniqueId uint32, characterId uint32, damage uint32) {
		e := &monsterDamageEvent{
			WorldId:     worldId,
			ChannelId:   channelId,
			MapId:       mapId,
			UniqueId:    uniqueId,
			CharacterId: characterId,
			Damage:      damage,
		}
		producer(CreateKey(int(uniqueId)), e)
	}
}
