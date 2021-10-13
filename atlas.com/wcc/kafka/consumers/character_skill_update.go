package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type characterSkillUpdateEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  int64  `json:"expiration"`
}

func CharacterSkillUpdateEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterSkillUpdateEvent{}
	}
}

func HandleCharacterSkillUpdateEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterSkillUpdateEvent); ok {
			if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, updateSkill(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func updateSkill(l logrus.FieldLogger, event *characterSkillUpdateEvent) session.Operator {
	return func(s *session.Model) {
		err := s.Announce(writer.WriteCharacterSkillUpdate(l)(event.SkillId, event.Level, event.MasterLevel, event.Expiration))
		if err != nil {
			l.WithError(err).Errorf("Unable to write skill update for character %d", event.CharacterId)
		}
	}
}
