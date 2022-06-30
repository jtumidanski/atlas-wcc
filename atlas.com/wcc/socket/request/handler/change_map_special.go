package handler

import (
	"atlas-wcc/character/properties"
	portal2 "atlas-wcc/portal"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
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

func ChangeMapSpecialHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readChangeMapSpecialRequest(r)
		c, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Cannot handle [ChangeMapSpecialRequest] because the acting character %d cannot be located.", s.CharacterId())
			return
		}

		portal, err := portal2.GetByName(l, span)(c.MapId(), p.StartWarp())
		if err != nil {
			l.WithError(err).Errorf("Cannot find portal %s in map %d in order to handle [ChangeMapSpecialRequest] for character %d", p.StartWarp(), c.MapId(), s.CharacterId())
			return
		}
		portal2.Enter(l, span)(worldId, channelId, c.MapId(), portal.Id(), s.CharacterId())
	}
}
