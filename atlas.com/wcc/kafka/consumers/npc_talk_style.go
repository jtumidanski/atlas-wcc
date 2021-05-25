package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
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
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, writeNpcTalkStyle(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeNpcTalkStyle(_ logrus.FieldLogger, event *npcTalkStyleCommand) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteNPCTalkStyle(event.NPCId, event.Message, event.Styles))
	}
}
