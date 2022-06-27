package skill

import (
	"atlas-wcc/kafka"
	"atlas-wcc/model"
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
		if _, err := session.GetByCharacterId(event.CharacterId); err != nil {
			return
		}

		session.ForSessionByCharacterId(event.CharacterId, showUpdate(l, event))
	}
}

func showUpdate(l logrus.FieldLogger, event updateEvent) model.Operator[session.Model] {
	return func(s session.Model) error {
		err := session.Announce(WriteCharacterSkillUpdate(l)(event.SkillId, event.Level, event.MasterLevel, event.Expiration))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to write skill update for character %d", event.CharacterId)
		}
		return err
	}
}
