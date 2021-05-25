package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*npcTalkNumCommand); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, writeNpcTalkNum(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeNpcTalkNum(_ logrus.FieldLogger, event *npcTalkNumCommand) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteNPCTalkNum(event.NPCId, event.Message, event.DefaultValue, event.MinimumValue, event.MaximumValue))
	}
}
