package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

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

func EmptyMagicAttackEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &magicAttackEvent{}
	}
}

func HandleMagicAttackEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*magicAttackEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			session.ForEachInMap(l)(event.WorldId, event.ChannelId, event.MapId, writeMagicAttack(l)(event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display, event.Charge))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeMagicAttack(l logrus.FieldLogger) func(characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte, charge int32) session.Operator {
	return func(characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte, charge int32) session.Operator {
		b := writer.WriteMagicAttack(l)(characterId, skill, skillLevel, stance, numberAttackedAndDamaged, damage, speed, direction, display, charge)
		return func(s *session.Model) {
			err := s.Announce(b)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		}
	}
}
