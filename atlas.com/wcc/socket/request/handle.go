package request

import (
	"atlas-wcc/account"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

type MessageValidator func(l logrus.FieldLogger, s *session.Model) bool

func NoOpValidator(_ logrus.FieldLogger, _ *session.Model) bool {
	return true
}

func LoggedInValidator(l logrus.FieldLogger, s *session.Model) bool {
	v := account.IsLoggedIn((*s).AccountId())
	if !v {
		l.Errorf("Attempting to process a request when the account %d is not logged in.", (*s).SessionId())
	}
	return v
}

type MessageHandler func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader)

func AdaptHandler(l logrus.FieldLogger, v MessageValidator, h MessageHandler) request.Handler {
	return func(sessionId uint32, r request.RequestReader) {
		s := session.GetRegistry().Get(sessionId)
		if s != nil {
			if v(l, s) {
				h(l, s, &r)
				s.UpdateLastRequest()
			}
		} else {
			l.Errorf("Unable to locate session %d", sessionId)
		}
	}
}
