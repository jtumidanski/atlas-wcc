package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type npcTalkStyleCommand struct {
	CharacterId uint32   `json:"characterId"`
	NPCId       uint32   `json:"npcId"`
	Message     string   `json:"message"`
	Type        string   `json:"type"`
	Speaker     string   `json:"speaker"`
	Styles      []uint32 `json:"styles"`
}

func EmptyNPCTalkStyleCommandCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &npcTalkStyleCommand{}
	}
}

func HandleNPCTalkStyleCommand() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*npcTalkStyleCommand); ok {
			if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, writeNpcTalkStyle(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeNpcTalkStyle(l logrus.FieldLogger, event *npcTalkStyleCommand) session.Operator {
	b := writer.WriteNPCTalkStyle(event.NPCId, event.Message, event.Styles)
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
