package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type CharacterSkillUpdateEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  int64  `json:"expiration"`
}

func CharacterSkillUpdateEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &CharacterSkillUpdateEvent{}
	}
}

func HandleCharacterSkillUpdateEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*CharacterSkillUpdateEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, updateSkill(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func updateSkill(_ logrus.FieldLogger, event *CharacterSkillUpdateEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteCharacterSkillUpdate(event.SkillId, event.Level, event.MasterLevel, event.Expiration))
	}
}
