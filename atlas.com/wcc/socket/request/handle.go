package request

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/registries"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

type SessionStateValidator func(l logrus.FieldLogger, s *mapleSession.MapleSession) bool

type SessionRequestHandler func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader)

func NoOpValidator() SessionStateValidator {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession) bool {
		return true
	}
}

func LoggedInValidator() SessionStateValidator {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession) bool {
		v := processors.IsLoggedIn((*s).AccountId())
		if !v {
			l.Errorf("Attempting to process a request when the account %d is not logged in.", (*s).SessionId())
		}
		return v
	}
}

func AdaptHandler(l logrus.FieldLogger, validator SessionStateValidator, handler SessionRequestHandler) request.Handler {
	return func(sessionId int, r request.RequestReader) {
		s := registries.GetSessionRegistry().Get(sessionId)
		if s != nil {
			if validator(l, &s) {
				handler(l, &s, &r)
				s.UpdateLastRequest()
			}
		} else {
			l.Errorf("Unable to locate session %d", sessionId)
		}
	}
}
