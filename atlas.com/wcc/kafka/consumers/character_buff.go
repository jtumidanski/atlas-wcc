package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

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

func CharacterBuffEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterBuffEvent{}
	}
}

func HandleCharacterBuffEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterBuffEvent); ok {
			session.ForSessionByCharacterId(event.CharacterId, showBuffEffect(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showBuffEffect(l logrus.FieldLogger, event *characterBuffEvent) session.Operator {
	b := writer.WriteShowBuff(l)(event.BuffId, event.Duration, makeBuffStats(event.Stats), event.Special)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func makeBuffStats(stats []stat) []writer.BuffStat {
	result := make([]writer.BuffStat, 0)
	for _, stat := range stats {
		result = append(result, makeBuffStat(stat))
	}
	return result
}

func makeBuffStat(s stat) writer.BuffStat {
	return writer.NewBuffStat(s.First, s.Mask, s.Amount)
}

type characterCancelBuffEvent struct {
	CharacterId uint32 `json:"characterId"`
	Stats       []stat `json:"stats"`
}

func CharacterCancelBuffEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterCancelBuffEvent{}
	}
}

func HandleCharacterCancelBuffEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterCancelBuffEvent); ok {
			session.ForSessionByCharacterId(event.CharacterId, cancelBuffEffect(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func cancelBuffEffect(l logrus.FieldLogger, event *characterCancelBuffEvent) session.Operator {
	return func(s *session.Model) {
		err := s.Announce(writer.WriteCancelBuff(l)(makeBuffStats(event.Stats)))
		if err != nil {
			l.WithError(err).Errorf("Unable to cancel buff for character %d", s.CharacterId())
		}
	}
}
