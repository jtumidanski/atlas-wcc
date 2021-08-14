package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type monsterMovementEvent struct {
	UniqueId      uint32      `json:"uniqueId"`
	ObserverId    uint32      `json:"observerId"`
	SkillPossible bool        `json:"skillPossible"`
	Skill         int8        `json:"skill"`
	SkillId       byte        `json:"skillId"`
	SkillLevel    byte        `json:"skillLevel"`
	Option        uint16      `json:"option"`
	StartX        int16       `json:"startX"`
	StartY        int16       `json:"startY"`
	EndX          int16       `json:"endX"`
	EndY          int16       `json:"endY"`
	Stance        byte        `json:"stance"`
	RawMovement   RawMovement `json:"rawMovement"`
}

type RawMovement []byte

func MonsterMovementEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &monsterMovementEvent{}
	}
}

func HandleMonsterMovementEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*monsterMovementEvent); ok {
			if actingSession := session.GetByCharacterId(event.ObserverId); actingSession == nil {
				return
			}

			session.ForEachOtherInMap(l)(wid, cid, event.ObserverId, moveMonster(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func moveMonster(l logrus.FieldLogger, event *monsterMovementEvent) session.Operator {
	b := writer.WriteMoveMonster(l)(event.UniqueId, event.SkillPossible, event.Skill, event.SkillId,
		event.SkillLevel, event.Option, event.StartX, event.StartY, event.RawMovement)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
