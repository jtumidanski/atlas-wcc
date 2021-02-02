package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	request2 "atlas-wcc/socket/request"
	"context"
	"github.com/jtumidanski/atlas-socket/request"
	"log"
)

const OpCharacterDistributeSp uint16 = 0x5A

type distributeSpRequest struct {
	skillId uint32
}

func (r distributeSpRequest) SkillId() uint32 {
	return r.skillId
}

func readDistributeSpRequest(reader *request.RequestReader) distributeSpRequest {
	reader.ReadUint32()
	skillId := reader.ReadUint32()
	return distributeSpRequest{skillId}
}

func DistributeSpHandler() request2.SessionRequestHandler {
	return func(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readDistributeSpRequest(r)
		producers.CharacterDistributeSp(l, context.Background()).Emit((*s).CharacterId(), p.SkillId())
	}
}
