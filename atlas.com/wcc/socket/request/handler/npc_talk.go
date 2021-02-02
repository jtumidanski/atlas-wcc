package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	request2 "atlas-wcc/socket/request"
	"atlas-wcc/socket/response/writer"
	"context"
	"github.com/jtumidanski/atlas-socket/request"
	"log"
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
	return func(l *log.Logger, s *mapleSession.MapleSession) bool {
		v := processors.IsLoggedIn((*s).AccountId())
		if !v {
			l.Printf("[ERROR] attempting to process a [HandleNPCTalkRequest] when the account %d is not logged in.", (*s).SessionId())
			(*s).Announce(writer.WriteEnableActions())
			return false
		}

		ca, err := processors.GetCharacterAttributesById((*s).CharacterId())
		if err != nil {
			l.Printf("[ERROR] unable to locate character %s speaking to npc.", (*s).CharacterId())
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
	return func(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readNPCTalkRequest(r)

		ca, err := processors.GetCharacterAttributesById((*s).CharacterId())
		if err != nil {
			l.Printf("[ERROR] unable to locate character %s speaking to npc.", (*s).CharacterId())
			return
		}

		npcs, err := processors.GetNPCsInMapByObjectId(ca.MapId(), p.ObjectId())
		if err != nil || len(npcs) != 1 {
			l.Printf("[ERROR] unable to locate npc %d in map %d.", p.ObjectId(), ca.MapId())
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

		if hasConversationScript(npc.Id()) {
			producers.NPCConversation(l, context.Background()).StartConversation((*s).WorldId(), (*s).ChannelId(), ca.MapId(), ca.Id(), npc.Id(), npc.ObjectId())
			return
		}

		//TODO deal with shops
	}
}

func hasConversationScript(npcId uint32) bool {
	return true
}

func handleGachapon(s *mapleSession.MapleSession) {

}

func handleDuey(s *mapleSession.MapleSession) {

}
