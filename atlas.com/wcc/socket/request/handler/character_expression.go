package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
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

func CharacterExpressionHandler() request2.SessionRequestHandler {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readCharacterExpressionRequest(r)
		producers.CharacterExpression(l)((*s).CharacterId(), p.Emote())
	}
}
