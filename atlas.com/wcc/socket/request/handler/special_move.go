package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpCodeSpecialMove uint16 = 0x5B

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

func HandleSpecialMove() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
		p := readSpecialMoveRequest(r)
		if event, ok := p.(monsterMagnetRequest); ok {
			data := make([]producers.MonsterMagnetData, 0)
			for _, d := range event.data {
				data = append(data, producers.MonsterMagnetData{MonsterId: d.monsterId, Success: d.success})
			}

			producers.ApplyMonsterMagnet(l)(s.CharacterId(), event.skillId, event.level, event.direction, data)
		} else if event, ok := p.(specialMoveRequest); ok {
			producers.ApplySkill(l)(s.CharacterId(), event.skillId, event.level, event.x, event.y)
		} else {
			l.Errorf("Received unexpected result from reading the special move request.")
		}
	}
}
