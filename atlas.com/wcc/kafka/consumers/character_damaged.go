package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type CharacterDamagedEvent struct {
	CharacterId     uint32 `json:"characterId"`
	MapId           uint32 `json:"mapId"`
	MonsterId       uint32 `json:"monsterId"`
	MonsterUniqueId uint32 `json:"monsterUniqueId"`
	SkillId         int8   `json:"skillId"`
	Damage          int32  `json:"damage"`
	Fake            uint32 `json:"fake"`
	Direction       int8   `json:"direction"`
	X               int16  `json:"x"`
	Y               int16  `json:"y"`
	PGMR            bool   `json:"pgmr"`
	PGMR1           byte   `json:"pgmr1"`
	PG              bool   `json:"pg"`
}

func CharacterDamagedEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &CharacterDamagedEvent{}
	}
}

func HandleCharacterDamagedEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*CharacterDamagedEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForEachSessionInMap(wid, cid, event.MapId, writeCharacterDamaged(l, *event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeCharacterDamaged(l logrus.FieldLogger, event CharacterDamagedEvent) session.SessionOperator {
	b := writer.WriteCharacterDamaged(event.SkillId, event.MonsterId, event.CharacterId, event.Damage, event.Fake, event.Direction, event.PGMR, event.PGMR1, event.PG, event.MonsterUniqueId, event.X, event.Y)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
