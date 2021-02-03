package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
)

type CharacterSkillUpdateEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  uint64 `json:"expiration"`
}

func CharacterSkillUpdateEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &CharacterSkillUpdateEvent{}
	}
}

func HandleCharacterSkillUpdateEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*CharacterSkillUpdateEvent); ok {
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}
			(*as).Announce(writer.WriteCharacterSkillUpdate(event.SkillId, event.Level, event.MasterLevel, event.Expiration))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterSkillUpdateEvent]")
		}
	}
}
