package request

import (
	"atlas-wcc/account"
	"atlas-wcc/model"
	"atlas-wcc/session"
	"atlas-wcc/tracing"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type MessageValidator func(l logrus.FieldLogger, span opentracing.Span) func(s session.Model) bool

func NoOpValidator(_ logrus.FieldLogger, _ opentracing.Span) func(_ session.Model) bool {
	return func(_ session.Model) bool {
		return true
	}
}

func LoggedInValidator(l logrus.FieldLogger, span opentracing.Span) func(s session.Model) bool {
	return func(s session.Model) bool {
		v := account.IsLoggedIn(l, span)(s.AccountId())
		if !v {
			l.Errorf("Attempting to process a request when the account %d is not logged in.", s.SessionId())
		}
		return v
	}
}

type MessageHandler func(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader)

func NoOpHandler(_ logrus.FieldLogger, _ opentracing.Span, _ byte, _ byte) func(_ session.Model, _ *request.RequestReader) {
	return func(_ session.Model, _ *request.RequestReader) {
	}
}

func AdaptHandler(l logrus.FieldLogger, worldId byte, channelId byte, name string, v MessageValidator, h MessageHandler) request.Handler {
	return func(sessionId uint32, r request.RequestReader) {
		sl, span := tracing.StartSpan(l, name)
		session.IfPresentById(sessionId, tryHandle(v(sl, span), h(sl, span, worldId, channelId), &r))
		span.Finish()
	}
}

func tryHandle(validator func(s session.Model) bool, handler func(_ session.Model, _ *request.RequestReader), reader *request.RequestReader) model.Operator[session.Model] {
	return func(s session.Model) error {
		if validator(s) {
			handler(s, reader)
		}
		s = session.UpdateLastRequest()(s.SessionId())
		return nil
	}
}
