package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/character/properties"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpGeneralChat uint16 = 0x31

type generalChatRequest struct {
	message string
	show    byte
}

func (r generalChatRequest) Message() string {
	return r.message
}

func (r generalChatRequest) Show() byte {
	return r.show
}

func readGeneralChatRequest(reader *request.RequestReader) generalChatRequest {
	message := reader.ReadAsciiString()
	show := reader.ReadByte()
	return generalChatRequest{message, show}
}

func GeneralChatHandler(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readGeneralChatRequest(r)
		ca, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Cannot handle [GeneralChatRequest] because the acting character %d cannot be located.", s.CharacterId())
			return
		}

		character.SendMapMessage(l, span)(s.WorldId(), s.ChannelId(), ca.MapId(), s.CharacterId(), p.Message(), ca.Gm(), p.Show() == 1)
	}
}
