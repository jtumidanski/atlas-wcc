package session

import (
	"atlas-wcc/model"
	"atlas-wcc/tracing"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net"
)

func Create(l logrus.FieldLogger) func(worldId byte, channelId byte) func(sessionId uint32, conn net.Conn) {
	return func(worldId byte, channelId byte) func(sessionId uint32, conn net.Conn) {
		return func(sessionId uint32, conn net.Conn) {
			l.Debugf("Creating session %d.", sessionId)

			s := NewSession(sessionId, conn)
			getRegistry().Add(s)

			err := s.WriteHello()
			if err != nil {
				l.WithError(err).Errorf("Unable to write hello packet to session %d.", sessionId)
			}
		}
	}
}

func Decrypt(_ logrus.FieldLogger) func(sessionId uint32, input []byte) []byte {
	return func(sessionId uint32, input []byte) []byte {
		s, err := GetById(sessionId)
		if err != nil {
			return input
		}

		if s.ReceiveAESOFB() == nil {
			return input
		}
		return s.ReceiveAESOFB().Decrypt(input, true, true)
	}
}

func DestroyAll(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte) {
	return func(worldId byte, channelId byte) {
		ForAll(Destroy(l, span)(worldId, channelId))
	}
}

func DestroyById(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte) func(sessionId uint32) {
	return func(worldId byte, channelId byte) func(sessionId uint32) {
		return func(sessionId uint32) {
			IfPresentById(sessionId, Destroy(l, span)(worldId, channelId))
		}
	}
}

func DestroyByIdWithSpan(l logrus.FieldLogger) func(worldId byte, channelId byte) func(sessionId uint32) {
	return func(worldId byte, channelId byte) func(sessionId uint32) {
		return func(sessionId uint32) {
			sl, span := tracing.StartSpan(l, "session_destroy")
			DestroyById(sl, span)(worldId, channelId)(sessionId)
			span.Finish()
		}
	}
}

func Destroy(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte) model.Operator[Model] {
	return func(worldId byte, channelId byte) model.Operator[Model] {
		return func(s Model) error {
			l.Debugf("Destroying session %d.", s.SessionId())

			getRegistry().Remove(s.SessionId())

			err := s.Disconnect()
			if err != nil {
				l.WithError(err).Errorf("Unable to issue disconnect to session %d.", s.SessionId())
				return err
			}

			Logout(l, span)(worldId, channelId, s.AccountId(), s.CharacterId())
			return nil
		}
	}
}
