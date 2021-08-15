package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/npc/conversation"
	"atlas-wcc/session"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpNpcTalkMore uint16 = 0x3C

type npcTalkMoreRequest struct {
	lastMessageType byte
	action          byte
	returnText      string
	selection       int32
}

func (r npcTalkMoreRequest) LastMessageType() byte {
	return r.lastMessageType
}

func (r npcTalkMoreRequest) Action() byte {
	return r.action
}

func (r npcTalkMoreRequest) ReturnText() string {
	return r.returnText
}

func (r npcTalkMoreRequest) Selection() int32 {
	return r.selection
}

func readNPCTalkMoreRequest(reader *request.RequestReader) npcTalkMoreRequest {
	lastMessageType := reader.ReadByte()
	action := reader.ReadByte()
	returnText := ""
	selection := int32(-1)

	if lastMessageType == 2 {
		if action != 0 {
			returnText = reader.ReadAsciiString()
		}
	} else {
		if len(reader.GetRestAsBytes()) >= 4 {
			selection = reader.ReadInt32()
		} else {
			selection = int32(reader.ReadByte())
		}
	}
	return npcTalkMoreRequest{lastMessageType, action, returnText, selection}
}

func HandleNPCTalkMoreRequest() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
		p := readNPCTalkMoreRequest(r)
		if p.LastMessageType() == 2 {
			if p.Action() != 0 {
				if questInProcess(s.CharacterId()) {
					continueQuestConversation(s.CharacterId(), p)
				} else {
					producers.SetReturnText(l)(s.CharacterId(), p.ReturnText())
					producers.ContinueConversation(l)(s.CharacterId(), p.Action(), p.LastMessageType(), -1)
				}
			} else if questInProcess(s.CharacterId()) {
				questDispose(s.CharacterId())
			} else {
				conversationDispose(s.CharacterId())
			}
		} else {
			if questInProcess(s.CharacterId()) {
				continueQuestConversation(s.CharacterId(), p)
			} else if conversation.InProgress(l)(s.CharacterId()) {
				producers.ContinueConversation(l)(s.CharacterId(), p.Action(), p.LastMessageType(), p.Selection())
			}
		}
	}
}

func conversationDispose(_ uint32) {

}

func questDispose(_ uint32) {

}

func questInProcess(_ uint32) bool {
	return false
}

func continueQuestConversation(_ uint32, _ npcTalkMoreRequest) {

}
