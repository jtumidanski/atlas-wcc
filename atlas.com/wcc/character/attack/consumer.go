package attack

import (
	"atlas-wcc/kafka"
	_map "atlas-wcc/map"
	"atlas-wcc/model"
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

		writer := writeCloseRangeAttack(l, event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display)
		_map.ForSessionsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, writer)
	}
}

func writeCloseRangeAttack(l logrus.FieldLogger, characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte) model.Operator[session.Model] {
	b := WriteCloseRangeAttack(l)(characterId, skill, skillLevel, stance, numberAttackedAndDamaged, damage, speed, direction, display)
	return func(s session.Model) error {
		err := session.Announce(b)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		return err
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

		_map.ForSessionsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, writeMagicAttack(l)(event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display, event.Charge))
	}
}

func writeMagicAttack(l logrus.FieldLogger) func(characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte, charge int32) model.Operator[session.Model] {
	return func(characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte, charge int32) model.Operator[session.Model] {
		b := WriteMagicAttack(l)(characterId, skill, skillLevel, stance, numberAttackedAndDamaged, damage, speed, direction, display, charge)
		return func(s session.Model) error {
			err := session.Announce(b)(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return err
		}
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

		_map.ForSessionsInMap(l, span)(event.WorldId, event.ChannelId, event.MapId, writeRangeAttack(l)(event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display, event.Projectile))
	}
}

func writeRangeAttack(l logrus.FieldLogger) func(characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte, projectile uint32) model.Operator[session.Model] {
	return func(characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte, projectile uint32) model.Operator[session.Model] {
		b := WriteRangeAttack(l)(characterId, skill, skillLevel, stance, numberAttackedAndDamaged, damage, speed, direction, display, projectile)
		return func(s session.Model) error {
			err := session.Announce(b)(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return err
		}
	}
}
