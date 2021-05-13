package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterBuffEvent struct {
	CharacterId uint32 `json:"characterId"`
	BuffId      uint32 `json:"id"`
	Duration    uint32 `json:"duration"`
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterBuffEvent); ok {
			processors.ForSessionByCharacterId(event.CharacterId, showBuffEffect(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showBuffEffect(_ logrus.FieldLogger, event *characterBuffEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteShowBuff(event.BuffId, event.Duration, makeBuffStats(event.Stats), event.Special))
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterCancelBuffEvent); ok {
			processors.ForSessionByCharacterId(event.CharacterId, cancelBuffEffect(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func cancelBuffEffect(_ logrus.FieldLogger, event *characterCancelBuffEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteCancelBuff(makeBuffStats(event.Stats)))
	}
}