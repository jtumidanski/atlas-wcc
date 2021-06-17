package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type closeRangeAttackEvent struct {
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

func EmptyCloseRangeAttackEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &closeRangeAttackEvent{}
	}
}

func HandleCloseRangeAttackEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*closeRangeAttackEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			session.ForEachInMap(event.WorldId, event.ChannelId, event.MapId, writeCloseRangeAttack(l, event.CharacterId, event.SkillId, event.SkillLevel, event.Stance, event.AttackedAndDamaged, event.Damage, event.Speed, event.Direction, event.Display))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeCloseRangeAttack(l logrus.FieldLogger, characterId uint32, skill uint32, skillLevel byte, stance byte, numberAttackedAndDamaged byte, damage map[uint32][]uint32, speed byte, direction byte, display byte) session.Operator {
	b := writer.WriteCloseRangeAttack(characterId, skill, skillLevel, stance, numberAttackedAndDamaged, damage, speed, direction, display)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
