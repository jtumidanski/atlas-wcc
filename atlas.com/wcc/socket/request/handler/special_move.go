package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpCodeSpecialMove uint16 = 0x5B
const SpecialMove = "special_move"

func HandleSpecialMoveProducer(l logrus.FieldLogger) Producer {
	return func() (uint16, request.Handler) {
		return OpCodeSpecialMove, SpanHandlerDecorator(l, SpecialMove, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(LoggedInValidator(l, span), HandleSpecialMove(l, span))
		})
	}
}

const (
	HeroMonsterMagnet       uint32 = 1121001
	PaladinMonsterMagnet    uint32 = 1221001
	DarkKnightMonsterMagnet uint32 = 1321001
	SuperGMHealPlusDispel   uint32 = 9101000
)

type monsterMagnetData struct {
	monsterId uint32
	success   uint8
}

type monsterMagnetRequest struct {
	skillId   uint32
	level     uint8
	direction int8
	data      []monsterMagnetData
}

type specialMoveRequest struct {
	skillId uint32
	level   uint8
	x       int16
	y       int16
}

func readSpecialMoveRequest(r *request.RequestReader) interface{} {
	r.ReadInt32()
	skillId := r.ReadUint32()
	level := r.ReadByte()

	if skillId == HeroMonsterMagnet || skillId == PaladinMonsterMagnet || skillId == DarkKnightMonsterMagnet {
		num := r.ReadUint32()
		md := make([]monsterMagnetData, 0)
		for i := uint32(0); i < num; i++ {
			monsterId := r.ReadUint32()
			success := r.ReadByte()
			md = append(md, monsterMagnetData{monsterId: monsterId, success: success})
		}
		direction := r.ReadInt8()
		return monsterMagnetRequest{
			skillId:   skillId,
			level:     level,
			direction: direction,
			data:      md,
		}
	} else if skillId == SuperGMHealPlusDispel {
		r.Skip(11)
	} else if skillId%10000000 == 1004 {
		r.ReadUint16()
	}

	x := int16(0)
	y := int16(0)
	if len(r.GetBuffer())-r.Position() == 5 {
		x = r.ReadInt16()
		y = r.ReadInt16()
	}
	return specialMoveRequest{
		skillId: skillId,
		level:   level,
		x:       x,
		y:       y,
	}
}

func HandleSpecialMove(l logrus.FieldLogger, span opentracing.Span) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readSpecialMoveRequest(r)
		if event, ok := p.(monsterMagnetRequest); ok {
			data := make([]character.MonsterMagnetData, 0)
			for _, d := range event.data {
				data = append(data, character.MonsterMagnetData{MonsterId: d.monsterId, Success: d.success})
			}

			character.ApplyMonsterMagnet(l, span)(s.CharacterId(), event.skillId, event.level, event.direction, data)
		} else if event, ok := p.(specialMoveRequest); ok {
			character.ApplySkill(l, span)(s.CharacterId(), event.skillId, event.level, event.x, event.y)
		} else {
			l.Errorf("Received unexpected result from reading the special move request.")
		}
	}
}
