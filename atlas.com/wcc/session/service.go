package session

import (
	"atlas-wcc/kafka/producers"
	"github.com/sirupsen/logrus"
	"net"
)

func Create(l logrus.FieldLogger, r *Registry) func(worldId byte, channelId byte) func(sessionId uint32, conn net.Conn) {
	return func(worldId byte, channelId byte) func(sessionId uint32, conn net.Conn) {
		return func(sessionId uint32, conn net.Conn) {
			l.Debugf("Creating session %d.", sessionId)
			s := NewSession(sessionId, conn)
			s.SetWorldId(worldId)
			s.SetChannelId(channelId)
			r.Add(s)
			s.WriteHello()
		}
	}
}

func Decrypt(_ logrus.FieldLogger, r *Registry) func(sessionId uint32, input []byte) []byte {
	return func(sessionId uint32, input []byte) []byte {
		s := r.Get(sessionId)
		if s == nil {
			return input
		}
		if s.ReceiveAESOFB() == nil {
			return input
		}
		return s.ReceiveAESOFB().Decrypt(input, true, true)
	}
}

func DestroyAll(l logrus.FieldLogger, r *Registry) {
	for _, s := range r.GetAll() {
		Destroy(l, r)(&s)
	}
}

func DestroyById(l logrus.FieldLogger, r *Registry) func(sessionId uint32) {
	return func(sessionId uint32) {
		s := r.Get(sessionId)
		if s == nil {
			return
		}
		Destroy(l, r)(s)
	}
}

func Destroy(l logrus.FieldLogger, r *Registry) func(session *Model) {
	return func(s *Model) {
		l.Debugf("Destroying session %d.", s.SessionId())
		r.Remove(s.SessionId())
		s.Disconnect()
		producers.Logout(l)(s.WorldId(), s.ChannelId(), s.AccountId(), s.CharacterId())
	}
}
