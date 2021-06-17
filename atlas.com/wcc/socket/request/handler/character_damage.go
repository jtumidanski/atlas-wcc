package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
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

func HandleCharacterDamageRequest() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
		p := readCharacterDamageRequest(r)
		producers.CharacterDamage(l)((*s).CharacterId(), p.MonsterIdFrom(), p.ObjectId(), p.DamageFrom(), p.Element(), p.Damage(), p.Direction())
	}
}
