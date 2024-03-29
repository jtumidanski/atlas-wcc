package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpCharacterHealOverTime uint16 = 0x59
const CharacterHealOverTime = "character_heal_over_time"

func HealOverTimeHandlerProducer(l logrus.FieldLogger) Producer {
	return func() (uint16, request.Handler) {
		return OpCharacterHealOverTime, SpanHandlerDecorator(l, CharacterHealOverTime, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(LoggedInValidator(l, span), HealOverTimeHandler(l, span))
		})
	}
}

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

func HealOverTimeHandler(l logrus.FieldLogger, span opentracing.Span) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readHealOverTimeRequest(r)

		character.AdjustHealth(l, span)(s.CharacterId(), p.HP())
		character.AdjustMana(l, span)(s.CharacterId(), p.MP())
	}
}
