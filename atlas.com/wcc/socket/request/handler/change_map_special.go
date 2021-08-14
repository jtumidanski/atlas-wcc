package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/kafka/producers"
	portal2 "atlas-wcc/portal"
	"atlas-wcc/session"
	request2 "atlas-wcc/socket/request"
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

func ChangeMapSpecialHandler() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
		p := readChangeMapSpecialRequest(r)
		c, err := character.GetCharacterAttributesById(l)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Cannot handle [ChangeMapSpecialRequest] because the acting character %d cannot be located.", s.CharacterId())
			return
		}

		portal, err := portal2.GetPortalByName(l)(c.MapId(), p.StartWarp())
		if err != nil {
			l.WithError(err).Errorf("Cannot find portal %s in map %d in order to handle [ChangeMapSpecialRequest] for character %d", p.StartWarp(), c.MapId(), s.CharacterId())
			return
		}
		producers.PortalEnter(l)(s.WorldId(), s.ChannelId(), c.MapId(), portal.Id(), s.CharacterId())
	}
}
