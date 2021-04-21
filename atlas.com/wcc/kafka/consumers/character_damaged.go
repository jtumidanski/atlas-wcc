package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type CharacterDamagedEvent struct {
	CharacterId     uint32 `json:"characterId"`
	MapId           uint32 `json:"mapId"`
	MonsterId       uint32 `json:"monsterId"`
	MonsterUniqueId uint32 `json:"monsterUniqueId"`
	SkillId         int8   `json:"skillId"`
	Damage          uint32 `json:"damage"`
	Fake            uint32 `json:"fake"`
	Direction       int8   `json:"direction"`
	X               int16  `json:"x"`
	Y               int16  `json:"y"`
	PGMR            bool   `json:"pgmr"`
	PGMR1           byte   `json:"pgmr1"`
	PG              bool   `json:"pg"`
}

func CharacterDamagedEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &CharacterDamagedEvent{}
	}
}

func HandleCharacterDamagedEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*CharacterDamagedEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachSessionInMap(wid, cid, event.MapId, writeCharacterDamaged(l, *event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeCharacterDamaged(_ logrus.FieldLogger, event CharacterDamagedEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteCharacterDamaged(event.SkillId, event.MonsterId, event.CharacterId, event.Damage,
			event.Fake, event.Direction, event.PGMR, event.PGMR1, event.PG, event.MonsterUniqueId, event.X, event.Y))
	}
}
