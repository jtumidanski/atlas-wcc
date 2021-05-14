package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	"atlas-wcc/npc/conversation"
	"atlas-wcc/npc/shop"
	"atlas-wcc/processors"
	request2 "atlas-wcc/socket/request"
	"atlas-wcc/socket/response/writer"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpNpcTalk uint16 = 0x3A

type npcTalkRequest struct {
	objectId uint32
}

func (r npcTalkRequest) ObjectId() uint32 {
	return r.objectId
}

func readNPCTalkRequest(reader *request.RequestReader) npcTalkRequest {
	return npcTalkRequest{reader.ReadUint32()}
}

func CharacterAliveValidator() request2.SessionStateValidator {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession) bool {
		v := processors.IsLoggedIn((*s).AccountId())
		if !v {
			l.Errorf("Attempting to process a [HandleNPCTalkRequest] when the account %d is not logged in.", (*s).SessionId())
			(*s).Announce(writer.WriteEnableActions())
			return false
		}

		ca, err := processors.GetCharacterAttributesById((*s).CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %s speaking to npc.", (*s).CharacterId())
			(*s).Announce(writer.WriteEnableActions())
			return false
		}

		if ca.Hp() > 0 {
			return true
		} else {
			(*s).Announce(writer.WriteEnableActions())
			return false
		}
	}
}

func HandleNPCTalkRequest() request2.SessionRequestHandler {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readNPCTalkRequest(r)

		ca, err := processors.GetCharacterAttributesById((*s).CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to locate character %s speaking to npc.", (*s).CharacterId())
			return
		}

		npcs, err := processors.GetNPCsInMapByObjectId(ca.MapId(), p.ObjectId())
		if err != nil || len(npcs) != 1 {
			l.WithError(err).Errorf("Unable to locate npc %d in map %d.", p.ObjectId(), ca.MapId())
			return
		}
		npc := npcs[0]

		if npc.Id() == 9010009 {
			handleDuey(s)
			return
		}

		if npc.Id() >= 9100100 && npc.Id() <= 9100200 {
			handleGachapon(s)
			return
		}

		if hasConversationScript(l)(npc.Id()) {
			producers.StartConversation(l)((*s).WorldId(), (*s).ChannelId(), ca.MapId(), ca.Id(), npc.Id(), npc.ObjectId())
			return
		}
		if hasShop(l)(npc.Id()) {
			ns, err := shop.GetShop(l)(npc.Id())
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve shop for npc %d.", npc.Id())
				return
			}
			err = (*s).Announce(writer.WriteGetNPCShop(ns))
			if err != nil {
				l.WithError(err).Errorf("Unable to write shop for npc %d to character %d.", npc.Id(), (*s).CharacterId())
			}
		}
	}
}

func hasShop(l logrus.FieldLogger) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		return shop.HasShop(l)(npcId)
	}
}

func hasConversationScript(l logrus.FieldLogger) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		return conversation.HasScript(l)(npcId)
	}
}

func handleGachapon(s *mapleSession.MapleSession) {

}

func handleDuey(s *mapleSession.MapleSession) {

}
