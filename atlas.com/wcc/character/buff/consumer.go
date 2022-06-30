package buff

import (
	"atlas-wcc/kafka"
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
		session.IfPresentByCharacterId(event.CharacterId, session.AnnounceOperator(WriteShowBuff(l)(event.BuffId, event.Duration, makeBuffStats(event.Stats), event.Special)))
	}
}

func makeBuffStats(stats []stat) []Stat {
	result := make([]Stat, 0)
	for _, s := range stats {
		result = append(result, makeBuffStat(s))
	}
	return result
}

func makeBuffStat(s stat) Stat {
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
		session.IfPresentByCharacterId(event.CharacterId, session.AnnounceOperator(WriteCancelBuff(l)(makeBuffStats(event.Stats))))
	}
}
