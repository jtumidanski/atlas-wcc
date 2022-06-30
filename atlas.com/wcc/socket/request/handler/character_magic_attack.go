package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/character/properties"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpCharacterMagicAttack uint16 = 0x2E
const CharacterMagicAttack = "character_magic_attack"

func CharacterMagicAttackHandlerProducer(l logrus.FieldLogger, worldId byte, channelId byte) Producer {
	return func() (uint16, request.Handler) {
		return OpCharacterMagicAttack, SpanHandlerDecorator(l, CharacterMagicAttack, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(LoggedInValidator(l, span), CharacterMagicAttackHandler(l, span, worldId, channelId))
		})
	}
}

func CharacterMagicAttackHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readAttackPacket(l, span, r, s.CharacterId(), false, true)

		catt, err := properties.GetById(l, span)(s.CharacterId())
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character attributes for character %d.", s.CharacterId())
			return
		}
		character.Attack(l, span)(worldId, channelId, catt.MapId(), s.CharacterId(), p.Skill(), p.SkillLevel(), p.NumberAttacked(), p.NumberDamaged(), p.NumberAttackedAndDamaged(), p.Stance(), p.Direction(), p.RangedDirection(), p.Charge(), p.Display(), p.Ranged(), p.Magic(), p.Speed(), p.AllDamage(), p.X(), p.Y())
	}
}
