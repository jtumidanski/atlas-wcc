package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
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
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*CharacterDamagedEvent); ok {
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}
			sl, err := getSessionsForThoseInMap(wid, cid, event.MapId)
			if err != nil {
				return
			}
			for _, s := range sl {
				s.Announce(writer.WriteCharacterDamaged(event.SkillId, event.MonsterId, event.CharacterId, event.Damage, event.Fake, event.Direction, event.PGMR, event.PGMR1, event.PG, event.MonsterUniqueId, event.X, event.Y))
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterDamagedEvent]")
		}
	}
}
