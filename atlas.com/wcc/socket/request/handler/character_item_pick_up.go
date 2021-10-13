package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpCharacterItemPickUp uint16 = 0xCA

type itemPickUpRequest struct {
	timestamp uint32
	x         int16
	y         int16
	objectId  uint32
}

func (r itemPickUpRequest) ObjectId() uint32 {
	return r.objectId
}

func readItemPickUpRequest(reader *request.RequestReader) itemPickUpRequest {
	timestamp := reader.ReadUint32()
	reader.ReadByte()
	x := reader.ReadInt16()
	y := reader.ReadInt16()
	objectId := reader.ReadUint32()
	return itemPickUpRequest{timestamp, x, y, objectId}
}

func ItemPickUpHandler(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readItemPickUpRequest(r)

		producers.CharacterReserveDrop(l, span)(s.CharacterId(), p.ObjectId())
	}
}
