package request

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/registries"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

type MessageValidator func(l logrus.FieldLogger, s *mapleSession.MapleSession) bool

func NoOpValidator(_ logrus.FieldLogger, _ *mapleSession.MapleSession) bool {
	return true
}

func LoggedInValidator(l logrus.FieldLogger, s *mapleSession.MapleSession) bool {
	v := processors.IsLoggedIn((*s).AccountId())
	if !v {
		l.Errorf("Attempting to process a request when the account %d is not logged in.", (*s).SessionId())
	}
	return v
}

type MessageHandler func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader)

func AdaptHandler(l logrus.FieldLogger, v MessageValidator, h MessageHandler) request.Handler {
	return func(sessionId uint32, r request.RequestReader) {
		s := registries.GetSessionRegistry().Get(sessionId)
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
