package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
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
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, updateSkill(l, event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterSkillUpdateEvent]")
		}
	}
}

func updateSkill(_ *log.Logger, event *CharacterSkillUpdateEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteCharacterSkillUpdate(event.SkillId, event.Level, event.MasterLevel, event.Expiration))
	}
}
