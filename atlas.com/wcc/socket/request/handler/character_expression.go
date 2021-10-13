package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpCharacterExpression uint16 = 0x33

type characterExpressionRequest struct {
	emote uint32
}

func (r characterExpressionRequest) Emote() uint32 {
	return r.emote
}

func readCharacterExpressionRequest(reader *request.RequestReader) characterExpressionRequest {
	emote := reader.ReadUint32()
	return characterExpressionRequest{emote}
}

func CharacterExpressionHandler(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readCharacterExpressionRequest(r)
		producers.CharacterExpression(l, span)(s.CharacterId(), p.Emote())
	}
}
