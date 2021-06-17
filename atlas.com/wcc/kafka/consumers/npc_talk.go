package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"fmt"
	"github.com/sirupsen/logrus"
)

type npcTalkEvent struct {
	CharacterId uint32 `json:"characterId"`
	NPCId       uint32 `json:"npcId"`
	Message     string `json:"message"`
	Type        string `json:"type"`
	Speaker     string `json:"speaker"`
}

func NPCTalkEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &npcTalkEvent{}
	}
}

func HandleNPCTalkEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*npcTalkEvent); ok {
			if actingSession := session.GetByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, writeNpcTalk(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeNpcTalk(l logrus.FieldLogger, event *npcTalkEvent) session.Operator {
	b := writer.WriteNPCTalk(event.NPCId, getNPCTalkType(event.Type), event.Message, getNPCTalkEnd(event.Type), getNPCTalkSpeaker(event.Speaker))
	return func(s *session.Model) {
		err := s.Announce(b)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func getNPCTalkSpeaker(speaker string) byte {
	switch speaker {
	case "NPC_LEFT":
		return 0
	case "NPC_RIGHT":
		return 1
	case "CHARACTER_LEFT":
		return 2
	case "CHARACTER_RIGHT":
		return 3
	}
	panic(fmt.Sprintf("unsupported npc talk speaker %s", speaker))
}

func getNPCTalkEnd(t string) []byte {
	switch t {
	case "NEXT":
		return []byte{00, 01}
	case "PREVIOUS":
		return []byte{01, 00}
	case "NEXT_PREVIOUS":
		return []byte{01, 01}
	case "OK":
		return []byte{00, 00}
	case "YES_NO":
		return []byte{}
	case "ACCEPT_DECLINE":
		return []byte{}
	case "SIMPLE":
		return []byte{}
	}
	panic(fmt.Sprintf("unsupported talk type %s", t))
}

func getNPCTalkType(t string) byte {
	switch t {
	case "NEXT":
		return 0
	case "PREVIOUS":
		return 0
	case "NEXT_PREVIOUS":
		return 0
	case "OK":
		return 0
	case "YES_NO":
		return 1
	case "ACCEPT_DECLINE":
		return 0x0C
	case "SIMPLE":
		return 4
	}
	panic(fmt.Sprintf("unsupported talk type %s", t))
}
