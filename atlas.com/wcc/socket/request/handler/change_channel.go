package handler

import (
	"atlas-wcc/channel"
	"atlas-wcc/character/properties"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpChangeChannel uint16 = 0x27

type changeChannelRequest struct {
	channelId byte
}

func (r changeChannelRequest) ChannelId() byte {
	return r.channelId
}

func readChangeChannelRequest(reader *request.RequestReader) changeChannelRequest {
	channelId := reader.ReadByte() + 1
	reader.ReadInt32()
	return changeChannelRequest{channelId}
}

func ChangeChannelHandler(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readChangeChannelRequest(r)
		if p.ChannelId() == s.ChannelId() {
			l.Errorf("Character %s trying to change to the same channel.", s.CharacterId())
			disconnect(l)(s)
		}

		//TODO further verification requests for ...
		// not being in cash shop
		// not being in mini game
		// not having a player shop open
		// not being in a mini dungeon

		char, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d changing channels.", s.CharacterId())
			disconnect(l)(s)
		}
		if char.Hp() <= 0 {
			l.Debugf("Character %d trying to change channel when dead.", s.CharacterId())
			return
		}

		ch, err := channel.GetForWorld(l, span)(s.WorldId(), p.ChannelId())
		if err != nil {
			l.WithError(err).Errorf("Cannot retrieve world %d channel %d information.", s.WorldId(), p.ChannelId())
			return
		}

		err = s.Announce(writer.WriteChangeChannel(l)(ch.IpAddress(), ch.Port()))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}

func disconnect(l logrus.FieldLogger) func(s *session.Model) {
	return func(s *session.Model) {
		err := s.Disconnect()
		if err != nil {
			l.WithError(err).Errorf("Unable to issue disconnect to session %d.", s.SessionId())
		}
	}
}
