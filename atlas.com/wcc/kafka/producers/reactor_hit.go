package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type hitReactorCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	MapId       uint32 `json:"map_id"`
	CharacterId uint32 `json:"character_id"`
	Id          uint32 `json:"id"`
	Stance      uint16 `json:"stance"`
	SkillId     uint32 `json:"skill_id"`
}

func HitReactor(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, id uint32, stance uint16, skillId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_HIT_REACTOR_COMMAND")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, id uint32, stance uint16, skillId uint32) {
		e := &hitReactorCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			MapId:       mapId,
			CharacterId: characterId,
			Id:          id,
			Stance:      stance,
			SkillId:     skillId,
		}
		producer(CreateKey(int(id)), e)
	}
}
