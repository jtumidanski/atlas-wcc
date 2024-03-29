package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/character/properties"
	"atlas-wcc/command"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpGeneralChat uint16 = 0x31
const GeneralChat = "general_chat"

func GeneralChatHandlerProducer(l logrus.FieldLogger, worldId byte, channelId byte) Producer {
	return func() (uint16, request.Handler) {
		return OpGeneralChat, SpanHandlerDecorator(l, GeneralChat, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(LoggedInValidator(l, span), GeneralChatHandler(l, span, worldId, channelId))
		})
	}
}

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

func GeneralChatHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readGeneralChatRequest(r)
		ca, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Cannot handle [GeneralChatRequest] because the acting character %d cannot be located.", s.CharacterId())
			return
		}

		e, found := command.Registry().Get(s, p.Message())
		if found {
			err = e(l, span)
			if err != nil {
				l.WithError(err).Errorf("Unable to execute command for character %d. Command=[%s]", s.CharacterId(), p.Message())
			}
			return
		}

		character.SendMapMessage(l, span)(worldId, channelId, ca.MapId(), s.CharacterId(), p.Message(), ca.Gm(), p.Show() == 1)
	}
}
