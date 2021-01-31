package services

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/registries"
	"github.com/jtumidanski/atlas-socket/session"
	"net"
)

type Service interface {
	session.Service
}

type mapleSessionService struct {
	r *registries.SessionRegistry
	worldId byte
	channelId byte
}

func NewMapleSessionService(wid byte, cid byte) Service {
	return &mapleSessionService{registries.GetSessionRegistry(), wid, cid}
}

func (s *mapleSessionService) Create(sessionId int, conn net.Conn) (session.Session, error) {
	ses := mapleSession.NewSession(sessionId, conn)
	s.r.Add(&ses)
	ses.SetWorldId(s.worldId)
	ses.SetChannelId(s.channelId)
	ses.WriteHello()
	return ses, nil
}

func (s *mapleSessionService) Get(sessionId int) session.Session {
	return s.r.Get(sessionId)
}

func (s *mapleSessionService) GetAll() []session.Session {
	ss := s.r.GetAll()
	b := make([]session.Session, len(ss))
	for i, v := range ss {
		b[i] = v.(session.Session)
	}
	return b
}

func (s *mapleSessionService) Destroy(sessionId int) {
	s.r.Remove(sessionId)
}
