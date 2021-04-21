package handler

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/socket/request"
	"atlas-wcc/socket/response/writer"
	request2 "github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpNpcAction uint16 = 0xC5

type npcAnimationRequest struct {
	first  uint32
	second byte
	third  byte
}

func (r npcAnimationRequest) First() uint32 {
	return r.first
}

func (r npcAnimationRequest) Second() byte {
	return r.second
}

func (r npcAnimationRequest) Third() byte {
	return r.third
}

type npcMoveRequest struct {
	movement []byte
}

func (r npcMoveRequest) Movement() []byte {
	return r.movement
}

func readNPCAction(reader *request2.RequestReader) interface{} {
	length := len(reader.GetRestAsBytes())
	if length == 6 {
		first := reader.ReadUint32()
		second := reader.ReadByte()
		third := reader.ReadByte()
		return &npcAnimationRequest{first, second, third}
	} else if length > 6 {
		bytes := reader.ReadBytes(length - 9)
		return &npcMoveRequest{bytes}
	}
	return nil
}

func HandleNPCAction() request.SessionRequestHandler {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request2.RequestReader) {
		p := readNPCAction(r)
		if val, ok := p.(*npcAnimationRequest); ok {
			(*s).Announce(writer.WriteNPCAnimation(val.First(), val.Second(), val.Third()))
		} else if val, ok := p.(*npcMoveRequest); ok {
			(*s).Announce(writer.WriteNPCMove(val.Movement()))
		} else {
			l.Warnf("Received a unhandled [NPCActionRequest]")
		}
	}
}
