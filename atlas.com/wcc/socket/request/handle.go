package request

import (
	"atlas-wcc/account"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type MessageValidator func(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model) bool

func NoOpValidator(_ logrus.FieldLogger, _ opentracing.Span) func(_ *session.Model) bool {
	return func(_ *session.Model) bool {
		return true
	}
}

func LoggedInValidator(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model) bool {
	return func(s *session.Model) bool {
		v := account.IsLoggedIn(l, span)(s.AccountId())
		if !v {
			l.Errorf("Attempting to process a request when the account %d is not logged in.", s.SessionId())
		}
		return v
	}
}

type MessageHandler func(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader)

func NoOpHandler(_ logrus.FieldLogger, _ opentracing.Span) func(_ *session.Model, _ *request.RequestReader) {
	return func(_ *session.Model, _ *request.RequestReader) {
	}
}

func AdaptHandler(l logrus.FieldLogger, name string, v MessageValidator, h MessageHandler) request.Handler {
	return func(sessionId uint32, r request.RequestReader) {
		s := session.Get().Get(sessionId)
		if s != nil {
			sp := opentracing.StartSpan(name)
			defer sp.Finish()

			if v(l, sp)(s) {
				h(l, sp)(s, &r)
				s.UpdateLastRequest()
			}
		} else {
			l.Errorf("Unable to locate session %d", sessionId)
		}
	}
}