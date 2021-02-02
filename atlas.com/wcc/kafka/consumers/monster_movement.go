package consumers

import (
	"atlas-wcc/rest/requests"
	"atlas-wcc/socket/response/writer"
	"log"
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

func MonsterMovementEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &monsterMovementEvent{}
	}
}

func HandleMonsterMovementEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*monsterMovementEvent); ok {
			m, err := requests.MonsterRegistry().GetById(event.UniqueId)
			if err != nil {
				l.Printf("[ERROR] unable to retrieve monster %d for MonsterMovementEvent", event.UniqueId)
				return
			}

			mapId := m.Data().Attributes.MapId
			sl, err := getSessionsForThoseInMap(wid, cid, mapId)
			if err != nil {
				l.Printf("[ERROR] unable to locate sessions for map %d-%d-%d.", wid, cid, mapId)
				return
			}
			for _, s := range sl {
				if s.CharacterId() != event.ObserverId {
					s.Announce(writer.WriteMoveMonster(event.UniqueId, event.SkillPossible, event.Skill, event.SkillId, event.SkillLevel, event.Option, event.StartX, event.StartY, event.RawMovement))
				}
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleMonsterMovementEvent]")
		}
	}
}
