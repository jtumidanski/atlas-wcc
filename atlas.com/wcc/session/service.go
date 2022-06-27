package session

import (
	"atlas-wcc/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net"
)

func Create(l logrus.FieldLogger, r *registry) func(worldId byte, channelId byte) func(sessionId uint32, conn net.Conn) {
	return func(worldId byte, channelId byte) func(sessionId uint32, conn net.Conn) {
		return func(sessionId uint32, conn net.Conn) {
			l.Debugf("Creating session %d.", sessionId)

			s := NewSession(sessionId, conn)
			r.Add(s)

			s = SetWorldId(worldId)(s.SessionId())
			s = SetChannelId(channelId)(s.SessionId())

			err := s.WriteHello()
			if err != nil {
				l.WithError(err).Errorf("Unable to write hello packet to session %d.", sessionId)
			}
		}
	}
}

func Decrypt(_ logrus.FieldLogger, r *registry) func(sessionId uint32, input []byte) []byte {
	return func(sessionId uint32, input []byte) []byte {
		s, ok := r.Get(sessionId)
		if !ok {
			return input
		}
		if s.ReceiveAESOFB() == nil {
			return input
		}
		return s.ReceiveAESOFB().Decrypt(input, true, true)
	}
}

func DestroyAll(l logrus.FieldLogger, span opentracing.Span, r *registry) {
	for _, s := range r.GetAll() {
		Destroy(l, span, r)(s)
	}
}

func DestroyById(l logrus.FieldLogger, span opentracing.Span, r *registry) func(sessionId uint32) {
	return func(sessionId uint32) {
		s, ok := r.Get(sessionId)
		if !ok {
			return
		}
		Destroy(l, span, r)(s)
	}
}

func DestroyByIdWithSpan(l logrus.FieldLogger, r *registry) func(sessionId uint32) {
	return func(sessionId uint32) {
		sl, span := tracing.StartSpan(l, "session_destroy")
		DestroyById(sl, span, r)(sessionId)
		span.Finish()
	}
}

func Destroy(l logrus.FieldLogger, span opentracing.Span, r *registry) func(session Model) {
	return func(s Model) {
		l.Debugf("Destroying session %d.", s.SessionId())

		r.Remove(s.SessionId())

		err := s.Disconnect()
		if err != nil {
			l.WithError(err).Errorf("Unable to issue disconnect to session %d.", s.SessionId())
		}

		Logout(l, span)(s.WorldId(), s.ChannelId(), s.AccountId(), s.CharacterId())
	}
}
