package handler

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	request2 "atlas-wcc/socket/request"
	"atlas-wcc/socket/response/writer"
	"github.com/jtumidanski/atlas-socket/request"
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

func ChangeChannelHandler() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readChangeChannelRequest(r)
		if p.ChannelId() == (*s).ChannelId() {
			l.Errorf("Character %s trying to change to the same channel.", (*s).CharacterId())
			(*s).Disconnect()
		}

		//TODO further verification requests for ...
		// not being in cash shop
		// not being in mini game
		// not having a player shop open
		// not being dead
		// not being in a mini dungeon

		channel, err := processors.GetChannelForWorld((*s).WorldId(), p.ChannelId())
		if err != nil {
			l.WithError(err).Errorf("Cannot retrieve world %d channel %d information.", (*s).WorldId(), p.ChannelId())
			return
		}

		(*s).Announce(writer.WriteChangeChannel(channel.IpAddress(), channel.Port()))
	}
}
