package handler

import (
	"atlas-wcc/party"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpPartyOperation uint16 = 0x7C

type createPartyPacket struct {
}

type leavePartyPacket struct {
}

type joinPartyPacket struct {
	partyId uint32
}

type invitePartyPacket struct {
	name string
}

type expelPartyPacket struct {
	characterId uint32
}

type changePartyLeaderPacket struct {
	newLeaderId uint32
}

func readPartyOperation(r *request.RequestReader) interface{} {
	op := r.ReadByte()
	switch op {
	case 1:
		return &createPartyPacket{}
	case 2:
		return &leavePartyPacket{}
	case 3:
		partyId := r.ReadUint32()
		return &joinPartyPacket{partyId: partyId}
	case 4:
		name := r.ReadAsciiString()
		return &invitePartyPacket{name: name}
	case 5:
		characterId := r.ReadUint32()
		return &expelPartyPacket{characterId: characterId}
	case 6:
		newLeaderId := r.ReadUint32()
		return &changePartyLeaderPacket{newLeaderId: newLeaderId}
	}
	return nil
}

func HandlePartyOperation(l logrus.FieldLogger, span opentracing.Span) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readPartyOperation(r)
		ok := false

		if _, ok = p.(*createPartyPacket); ok {
			party.Create(l, span)(s.WorldId(), s.ChannelId(), s.CharacterId())
			return
		} else if _, ok = p.(*leavePartyPacket); ok {
			party.Leave(l, span)(s.WorldId(), s.ChannelId(), s.CharacterId())
			return
		}

		l.Warnf("Received a unhandled operation %v.", &p)
	}
}
