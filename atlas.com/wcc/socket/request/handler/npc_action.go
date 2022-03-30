package handler

import (
	"atlas-wcc/npc"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
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

func readNPCAction(reader *request.RequestReader) interface{} {
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

func HandleNPCAction(l logrus.FieldLogger, _ opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readNPCAction(r)
		if val, ok := p.(*npcAnimationRequest); ok {
			err := s.Announce(npc.WriteNPCAnimation(l)(val.ObjectId(), val.Second(), val.Third()))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		} else if val, ok := p.(*npcMoveRequest); ok {
			err := s.Announce(npc.WriteNPCMove(l)(val.Movement()))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		} else {
			l.Warnf("Received a unhandled [NPCActionRequest]")
		}
	}
}
