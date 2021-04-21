package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	request2 "atlas-wcc/socket/request"
	"context"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpChangeMapSpecial uint16 = 0x64

type changeMapSpecialRequest struct {
	startWarp string
}

func (c *changeMapSpecialRequest) StartWarp() string {
	return c.startWarp
}

func readChangeMapSpecialRequest(reader *request.RequestReader) changeMapSpecialRequest {
	reader.ReadByte()
	sw := reader.ReadAsciiString()
	reader.ReadUint16()
	return changeMapSpecialRequest{sw}
}

func ChangeMapSpecialHandler() request2.SessionRequestHandler {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readChangeMapSpecialRequest(r)
		c, err := processors.GetCharacterAttributesById((*s).CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Cannot handle [ChangeMapSpecialRequest] because the acting character %d cannot be located.", (*s).CharacterId())
			return
		}

		portal, err := processors.GetPortalByName(c.MapId(), p.StartWarp())
		if err != nil {
			l.WithError(err).Errorf("Cannot find portal %s in map %d in order to handle [ChangeMapSpecialRequest] for character %d", p.StartWarp(), c.MapId(), (*s).CharacterId())
			return
		}
		producers.PortalEnter(l, context.Background()).Enter((*s).WorldId(), (*s).ChannelId(), c.MapId(), portal.Id(), (*s).CharacterId())
	}
}
