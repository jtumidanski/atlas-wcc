package reactor

import (
	"atlas-wcc/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type hitCommand struct {
	WorldId     byte   `json:"world_id"`
	ChannelId   byte   `json:"channel_id"`
	MapId       uint32 `json:"map_id"`
	CharacterId uint32 `json:"character_id"`
	Id          uint32 `json:"id"`
	Stance      uint16 `json:"stance"`
	SkillId     uint32 `json:"skill_id"`
}

func Hit(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, id uint32, stance uint16, skillId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_HIT_REACTOR_COMMAND")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, id uint32, stance uint16, skillId uint32) {
		e := &hitCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			MapId:       mapId,
			CharacterId: characterId,
			Id:          id,
			Stance:      stance,
			SkillId:     skillId,
		}
		producer(kafka.CreateKey(int(id)), e)
	}
}
