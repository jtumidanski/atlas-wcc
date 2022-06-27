package buff

import (
	"atlas-wcc/kafka"
	"atlas-wcc/model"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	consumerNameCharacterBuff       = "character_buff_event"
	consumerNameCharacterCancelBuff = "character_cancel_buff_event"
	topicTokenCharacterBuff         = "TOPIC_CHARACTER_BUFF"
	topicTokenCancelCharacterBuff   = "TOPIC_CHARACTER_CANCEL_BUFF"
)

func CharacterBuffConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterBuffEvent](consumerNameCharacterBuff, topicTokenCharacterBuff, groupId, handleCharacterBuff(wid, cid))
	}
}

type characterBuffEvent struct {
	CharacterId uint32 `json:"characterId"`
	BuffId      uint32 `json:"id"`
	Duration    int32  `json:"duration"`
	Stats       []stat `json:"stats"`
	Special     bool   `json:"special"`
}

type stat struct {
	First  bool   `json:"first"`
	Mask   uint64 `json:"mask"`
	Amount uint16 `json:"amount"`
}

func handleCharacterBuff(_ byte, _ byte) kafka.HandlerFunc[characterBuffEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event characterBuffEvent) {
		session.ForSessionByCharacterId(event.CharacterId, showBuffEffect(l, event))
	}
}

func showBuffEffect(l logrus.FieldLogger, event characterBuffEvent) model.Operator[session.Model] {
	b := WriteShowBuff(l)(event.BuffId, event.Duration, makeBuffStats(event.Stats), event.Special)
	return func(s session.Model) error {
		err := session.Announce(b)(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		return err
	}
}

func makeBuffStats(stats []stat) []BuffStat {
	result := make([]BuffStat, 0)
	for _, s := range stats {
		result = append(result, makeBuffStat(s))
	}
	return result
}

func makeBuffStat(s stat) BuffStat {
	return NewBuffStat(s.First, s.Mask, s.Amount)
}

func CharacterCancelBuffConsumer(wid byte, cid byte) func(groupId string) kafka.ConsumerConfig {
	return func(groupId string) kafka.ConsumerConfig {
		return kafka.NewConsumerConfig[characterCancelBuffEvent](consumerNameCharacterCancelBuff, topicTokenCancelCharacterBuff, groupId, handleCharacterCancelBuff(wid, cid))
	}
}

type characterCancelBuffEvent struct {
	CharacterId uint32 `json:"characterId"`
	Stats       []stat `json:"stats"`
}

func handleCharacterCancelBuff(_ byte, _ byte) kafka.HandlerFunc[characterCancelBuffEvent] {
	return func(l logrus.FieldLogger, span opentracing.Span, event characterCancelBuffEvent) {
		session.ForSessionByCharacterId(event.CharacterId, cancelBuffEffect(l, event))
	}
}

func cancelBuffEffect(l logrus.FieldLogger, event characterCancelBuffEvent) model.Operator[session.Model] {
	return func(s session.Model) error {
		err := session.Announce(WriteCancelBuff(l)(makeBuffStats(event.Stats)))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to cancel buff for character %d", s.CharacterId())
		}
		return err
	}
}
