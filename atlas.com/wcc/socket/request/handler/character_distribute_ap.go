package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpCharacterDistributeAp uint16 = 0x57

type distributeApRequest struct {
	number uint32
}

func (r distributeApRequest) Number() uint32 {
	return r.number
}

func readDistributeApRequest(reader *request.RequestReader) distributeApRequest {
	reader.ReadUint32()
	number := reader.ReadUint32()
	return distributeApRequest{number}
}

func DistributeApHandler(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readDistributeApRequest(r)

		attributeType := getType(p.Number())
		producers.CharacterDistributeAp(l, span)(s.CharacterId(), attributeType)
	}
}

func getType(number uint32) string {
	switch number {
	case 64:
		return "STRENGTH"
	case 128:
		return "DEXTERITY"
	case 256:
		return "INTELLIGENCE"
	case 512:
		return "LUCK"
	case 2048:
		return "HP"
	case 8192:
		return "MP"
	default:
		panic("invalid distribute ap value")
	}
}
