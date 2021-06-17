package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpCharacterHealOverTime uint16 = 0x59

type healOverTimeRequest struct {
	hp int16
	mp int16
}

func (r healOverTimeRequest) HP() int16 {
	return r.hp
}

func (r healOverTimeRequest) MP() int16 {
	return r.mp
}

func readHealOverTimeRequest(reader *request.RequestReader) healOverTimeRequest {
	reader.Skip(8)
	hp := reader.ReadInt16()
	mp := reader.ReadInt16()
	return healOverTimeRequest{hp, mp}
}

func HealOverTimeHandler() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
		p := readHealOverTimeRequest(r)

		producers.CharacterAdjustHealth(l)((*s).CharacterId(), p.HP())
		producers.CharacterAdjustMana(l)((*s).CharacterId(), p.MP())
	}
}
