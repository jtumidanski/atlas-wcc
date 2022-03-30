package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
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

func DistributeSpHandler(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readDistributeSpRequest(r)
		character.DistributeSp(l, span)(s.CharacterId(), p.SkillId())
	}
}
