package handler

import (
	"atlas-wcc/character/properties"
	"atlas-wcc/kafka/producers"
	portal2 "atlas-wcc/portal"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpChangeMap uint16 = 0x26

type changeMapRequest struct {
	cashShop  bool
	fromDying int8
	targetId  int32
	startWarp string
	wheel     bool
}

func (r changeMapRequest) CashShop() bool {
	return r.cashShop
}

func (r changeMapRequest) StartWarp() string {
	return r.startWarp
}

func readChangeMapRequest(reader *request.RequestReader) changeMapRequest {
	cs := len(reader.String()) == 0
	fromDying := int8(-1)
	targetId := int32(-1)
	startWarp := ""
	wheel := false

	if !cs {
		fromDying = reader.ReadInt8()
		targetId = reader.ReadInt32()
		startWarp = reader.ReadAsciiString()
		reader.ReadByte()
		wheel = reader.ReadInt16() > 0
	}
	return changeMapRequest{cs, fromDying, targetId, startWarp, wheel}
}

func ChangeMapHandler(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readChangeMapRequest(r)
		if p.CashShop() {

		} else {
			ca, err := properties.GetById(l, span)(s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Cannot handle [ChangeMapRequest] because the acting character %d cannot be located.", s.CharacterId())
				return
			}

			portal, err := portal2.GetByName(l, span)(ca.MapId(), p.StartWarp())
			if err != nil {
				l.WithError(err).Errorf("Cannot find portal %s in map %d in order to handle [ChangeMapRequest] for character %d", p.StartWarp(), ca.MapId(), s.CharacterId())
				return
			}
			producers.PortalEnter(l, span)(s.WorldId(), s.ChannelId(), ca.MapId(), portal.Id(), s.CharacterId())
		}
	}
}
