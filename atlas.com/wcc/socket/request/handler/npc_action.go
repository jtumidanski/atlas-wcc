package handler

import (
	"atlas-wcc/session"
	"atlas-wcc/socket/request"
	"atlas-wcc/socket/response/writer"
	request2 "github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpNpcAction uint16 = 0xC5

type npcAnimationRequest struct {
	objectId uint32
	second   byte
	third    byte
}

func (r npcAnimationRequest) ObjectId() uint32 {
	return r.objectId
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

func HandleNPCAction() request.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request2.RequestReader) {
		p := readNPCAction(r)
		if val, ok := p.(*npcAnimationRequest); ok {
			err := s.Announce(writer.WriteNPCAnimation(l)(val.ObjectId(), val.Second(), val.Third()))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		} else if val, ok := p.(*npcMoveRequest); ok {
			err := s.Announce(writer.WriteNPCMove(l)(val.Movement()))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		} else {
			l.Warnf("Received a unhandled [NPCActionRequest]")
		}
	}
}
