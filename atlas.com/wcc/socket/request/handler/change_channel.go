package handler

import (
	"atlas-wcc/channel"
	"atlas-wcc/character/properties"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpChangeChannel uint16 = 0x27
const ChangeChannel = "change_channel"

func ChangeChannelHandlerProducer(l logrus.FieldLogger, worldId byte, channelId byte) Producer {
	return func() (uint16, request.Handler) {
		return OpChangeChannel, SpanHandlerDecorator(l, ChangeChannel, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(LoggedInValidator(l, span), ChangeChannelHandler(l, span, worldId, channelId))
		})
	}
}

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

func ChangeChannelHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readChangeChannelRequest(r)
		if p.ChannelId() == channelId {
			l.Errorf("Character %s trying to change to the same channel. Disconnect them.", s.CharacterId())
			_ = session.Destroy(l, span)(worldId, channelId)(s)
		}

		//TODO further verification requests for ...
		// not being in cash shop
		// not being in mini game
		// not having a player shop open
		// not being in a mini dungeon

		char, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d changing channels.", s.CharacterId())
			_ = session.Destroy(l, span)(worldId, channelId)(s)
		}
		if char.Hp() <= 0 {
			l.Debugf("Character %d trying to change channel when dead.", s.CharacterId())
			return
		}

		ch, err := channel.GetByWorldId(l, span)(worldId, p.ChannelId())
		if err != nil {
			l.WithError(err).Errorf("Cannot retrieve world %d channel %d information.", worldId, p.ChannelId())
			return
		}

		err = session.Announce(s, channel.WriteChangeChannel(l)(ch.IpAddress(), ch.Port()))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
	}
}
