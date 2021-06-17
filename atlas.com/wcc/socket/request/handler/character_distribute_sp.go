package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
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

func DistributeSpHandler() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readDistributeSpRequest(r)
		producers.CharacterDistributeSp(l)((*s).CharacterId(), p.SkillId())
	}
}
