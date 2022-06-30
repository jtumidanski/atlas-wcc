package handler

import (
	"atlas-wcc/account"
	"atlas-wcc/session"
	"atlas-wcc/tracing"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type MessageValidator func(s session.Model) bool

func NoOpValidator(_ session.Model) bool {
	return true
}

func LoggedInValidator(l logrus.FieldLogger, span opentracing.Span) MessageValidator {
	return func(s session.Model) bool {
		v := account.IsLoggedIn(l, span)(s.AccountId())
		if !v {
			l.Errorf("Attempting to process a request when the account %d is not logged in.", s.SessionId())
		}
		return v
	}

}

type MessageHandler func(s session.Model, r *request.RequestReader)

func NoOpHandler(_ session.Model, _ *request.RequestReader) {
}

type SpanAwareHandler func(l logrus.FieldLogger, span opentracing.Span) request.Handler

func SpanHandlerDecorator(l logrus.FieldLogger, name string, handler SpanAwareHandler) request.Handler {
	return func(sessionId uint32, reader request.RequestReader) {
		sl, span := tracing.StartSpan(l, name)
		handler(sl, span)(sessionId, reader)
	}
}

func ValidatorHandler(v MessageValidator, h MessageHandler) request.Handler {
	return func(sessionId uint32, reader request.RequestReader) {
		session.IfPresentById(sessionId, func(s session.Model) error {
			if v(s) {
				h(s, &reader)
			}
			s = session.UpdateLastRequest()(s.SessionId())
			return nil
		})
	}
}
