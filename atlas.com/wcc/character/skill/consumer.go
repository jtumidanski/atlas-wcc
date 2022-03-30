package skill

import (
	"atlas-wcc/kafka"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameUpdate = "character_skill_update_event"
	topicTokenUpdate   = "TOPIC_CHARACTER_SKILL_UPDATE_EVENT"
)

func UpdateConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[updateEvent](consumerNameUpdate, topicTokenUpdate, groupId, handleUpdate(wid, cid))
	}
}

type updateEvent struct {
	CharacterId uint32 `json:"characterId"`
	SkillId     uint32 `json:"skillId"`
	Level       uint32 `json:"level"`
	MasterLevel uint32 `json:"masterLevel"`
	Expiration  int64  `json:"expiration"`
}

func handleUpdate(_ byte, _ byte) kafka.HandlerFunc[updateEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event updateEvent) {
		if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, showUpdate(l, event))
	}
}

func showUpdate(l logrus.FieldLogger, event updateEvent) session.Operator {
	return func(s *session.Model) {
		err := s.Announce(WriteCharacterSkillUpdate(l)(event.SkillId, event.Level, event.MasterLevel, event.Expiration))
		if err != nil {
			l.WithError(err).Errorf("Unable to write skill update for character %d", event.CharacterId)
		}
	}
}
