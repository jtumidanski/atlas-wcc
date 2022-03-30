package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpCharacterDamage uint16 = 0x30

type characterDamageRequest struct {
	damageFrom    int8
	element       byte
	damage        int32
	monsterIdFrom uint32
	objectId      uint32
	direction     int8
}

func (r characterDamageRequest) MonsterIdFrom() uint32 {
	return r.monsterIdFrom
}

func (r characterDamageRequest) ObjectId() uint32 {
	return r.objectId
}

func (r characterDamageRequest) DamageFrom() int8 {
	return r.damageFrom
}

func (r characterDamageRequest) Element() byte {
	return r.element
}

func (r characterDamageRequest) Damage() int32 {
	return r.damage
}

func (r characterDamageRequest) Direction() int8 {
	return r.direction
}

func readCharacterDamageRequest(reader *request.RequestReader) characterDamageRequest {
	reader.ReadUint32()
	damageFrom := reader.ReadInt8()
	element := reader.ReadByte()
	damage := reader.ReadInt32()
	monsterIdFrom := uint32(0)
	oid := uint32(0)
	if damageFrom != -3 && damageFrom != -4 {
		monsterIdFrom = reader.ReadUint32()
		oid = reader.ReadUint32()
	}
	direction := reader.ReadInt8()
	return characterDamageRequest{
		damageFrom:    damageFrom,
		element:       element,
		damage:        damage,
		monsterIdFrom: monsterIdFrom,
		objectId:      oid,
		direction:     direction,
	}
}

func HandleCharacterDamageRequest(l logrus.FieldLogger, span opentracing.Span) func(s *session.Model, r *request.RequestReader) {
	return func(s *session.Model, r *request.RequestReader) {
		p := readCharacterDamageRequest(r)
		character.Damage(l, span)(s.CharacterId(), p.MonsterIdFrom(), p.ObjectId(), p.DamageFrom(), p.Element(), p.Damage(), p.Direction())
	}
}
