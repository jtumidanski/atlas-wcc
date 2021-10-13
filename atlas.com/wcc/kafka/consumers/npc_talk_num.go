package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type npcTalkNumCommand struct {
	CharacterId  uint32 `json:"characterId"`
	NPCId        uint32 `json:"npcId"`
	Message      string `json:"message"`
	Type         string `json:"type"`
	Speaker      string `json:"speaker"`
	DefaultValue int32  `json:"defaultValue"`
	MinimumValue int32  `json:"minimumValue"`
	MaximumValue int32  `json:"maximumValue"`
}

func EmptyNPCTalkNumCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &npcTalkNumCommand{}
	}
}

func HandleNPCTalkNumCommand() ChannelEventProcessor {
	return func(l logrus.FieldLogger, span opentracing.Span, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*npcTalkNumCommand); ok {
			if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, writeNpcTalkNum(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeNpcTalkNum(l logrus.FieldLogger, event *npcTalkNumCommand) session.Operator {
	b := writer.WriteNPCTalkNum(l)(event.NPCId, event.Message, event.DefaultValue, event.MinimumValue, event.MaximumValue)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
