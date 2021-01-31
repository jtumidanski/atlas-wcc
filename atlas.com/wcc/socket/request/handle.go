package request

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/registries"
	"github.com/jtumidanski/atlas-socket/request"
	"log"
)

type HandlerSupplier struct {
	l *log.Logger
}

type MapleHandler interface {
	IsValid(l *log.Logger, s *mapleSession.MapleSession) bool

	HandleRequest(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader)
}

func AdaptHandler(l *log.Logger, h MapleHandler) func(int, request.RequestReader) {
	return func(sessionId int, r request.RequestReader) {
		s := registries.GetSessionRegistry().Get(sessionId)
		if s != nil {
			if h.IsValid(l, &s) {
				h.HandleRequest(l, &s, &r)
				s.UpdateLastRequest()
			}
		} else {
			l.Printf("[ERROR] unable to locate session %d", sessionId)
		}
	}
}
