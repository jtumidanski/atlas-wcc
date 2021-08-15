package handler

import (
	"atlas-wcc/character/properties"
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	"atlas-wcc/socket/request"
	request2 "github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpReactorHit uint16 = 0xCD

type reactorHitRequest struct {
	id                uint32
	characterPosition int32
	stance            uint16
	skillId           uint32
}

func (r reactorHitRequest) Id() uint32 {
	return r.id
}

func (r reactorHitRequest) Stance() uint16 {
	return r.stance
}

func (r reactorHitRequest) SkillId() uint32 {
	return r.skillId
}

func readReactorHit(r *request2.RequestReader) interface{} {
	id := r.ReadUint32()
	characterPosition := r.ReadInt32()
	stance := r.ReadUint16()
	r.Skip(4)
	skillId := r.ReadUint32()
	return &reactorHitRequest{
		id:                id,
		characterPosition: characterPosition,
		stance:            stance,
		skillId:           skillId,
	}
}

func HandleReactorHit() request.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request2.RequestReader) {
		p := readReactorHit(r)
		if val, ok := p.(*reactorHitRequest); ok {
			c, err := properties.GetById(l)(s.CharacterId())
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve character %d for session %d.", s.CharacterId(), s.SetCharacterId)
				return
			}

			producers.HitReactor(l)(
				s.WorldId(),
				s.ChannelId(),
				c.MapId(),
				s.CharacterId(),
				val.Id(),
				val.Stance(),
				val.SkillId(),
			)
		}
	}
}
