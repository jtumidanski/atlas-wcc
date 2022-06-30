package attack

import (
	"atlas-wcc/kafka"
	_map "atlas-wcc/map"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameCloseRangeAttack = "close_range_attack_event"
	consumerNameRangeAttack      = "range_attack_event"
	consumerNameMagicAttack      = "magic_attack_event"
	topicTokenCloseRangeAttack   = "TOPIC_CLOSE_RANGE_ATTACK_EVENT"
	topicTokenRangeAttack        = "TOPIC_RANGE_ATTACK_EVENT"
	topicTokenMagicAttack        = "TOPIC_MAGIC_ATTACK_EVENT"
)

func CloseRangeConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[closeRangeEvent](consumerNameCloseRangeAttack, topicTokenCloseRangeAttack, groupId, handleCloseRange(wid, cid))
	}
}

type closeRangeEvent struct {
	WorldId            byte                `json:"worldId"`
	ChannelId          byte                `json:"channelId"`
	MapId              uint32              `json:"mapId"`
	CharacterId        uint32              `json:"characterId"`
	SkillId            uint32              `json:"skillId"`
	SkillLevel         byte                `json:"skillLevel"`
	AttackedAndDamaged byte                `json:"attackedAndDamaged"`
	Display            byte                `json:"display"`
	Direction          byte                `json:"direction"`
	Stance             byte                `json:"stance"`
	Speed              byte                `json:"speed"`
	Damage             map[uint32][]uint32 `json:"damage"`
}

func handleCloseRange(wid byte, cid byte) kafka.HandlerFunc[closeRangeEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event closeRangeEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		_map.ForSessionsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, session.Announce(WriteCloseRangeAttack(l)(event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display)))
	}
}

func MagicAttackConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[magicAttackEvent](consumerNameMagicAttack, topicTokenMagicAttack, groupId, handleMagicAttack(wid, cid))
	}
}

type magicAttackEvent struct {
	WorldId            byte                `json:"worldId"`
	ChannelId          byte                `json:"channelId"`
	MapId              uint32              `json:"mapId"`
	CharacterId        uint32              `json:"characterId"`
	SkillId            uint32              `json:"skillId"`
	SkillLevel         byte                `json:"skillLevel"`
	Stance             byte                `json:"stance"`
	AttackedAndDamaged byte                `json:"attackedAndDamaged"`
	Damage             map[uint32][]uint32 `json:"damage"`
	Speed              byte                `json:"speed"`
	Direction          byte                `json:"direction"`
	Display            byte                `json:"display"`
	Charge             int32               `json:"charge"`
}

func handleMagicAttack(wid byte, cid byte) kafka.HandlerFunc[magicAttackEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event magicAttackEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		_map.ForSessionsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, session.Announce(WriteMagicAttack(l)(event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display, event.Charge)))
	}
}

func RangeAttackConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[rangeAttackEvent](consumerNameRangeAttack, topicTokenRangeAttack, groupId, handleRangeAttack(wid, cid))
	}
}

type rangeAttackEvent struct {
	WorldId            byte                `json:"worldId"`
	ChannelId          byte                `json:"channelId"`
	MapId              uint32              `json:"mapId"`
	CharacterId        uint32              `json:"characterId"`
	SkillId            uint32              `json:"skillId"`
	SkillLevel         byte                `json:"skillLevel"`
	Stance             byte                `json:"stance"`
	AttackedAndDamaged byte                `json:"attackedAndDamaged"`
	Projectile         uint32              `json:"projectile"`
	Damage             map[uint32][]uint32 `json:"damage"`
	Speed              byte                `json:"speed"`
	Direction          byte                `json:"direction"`
	Display            byte                `json:"display"`
}

func handleRangeAttack(wid byte, cid byte) kafka.HandlerFunc[rangeAttackEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event rangeAttackEvent) {
		if wid != event.WorldId || cid != event.ChannelId {
			return
		}

		_map.ForSessionsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, session.Announce(WriteRangeAttack(l)(event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display, event.Projectile)))
	}
}
