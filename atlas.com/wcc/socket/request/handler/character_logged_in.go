package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/character/keymap"
	"atlas-wcc/map"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpCharacterLoggedIn uint16 = 0x14
const NameCharacterLoggedIn = "character_logged_in"

func CharacterLoggedInHandlerProducer(l logrus.FieldLogger, worldId byte, channelId byte) Producer {
	return func() (uint16, request.Handler) {
		return OpCharacterLoggedIn, SpanHandlerDecorator(l, NameCharacterLoggedIn, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(NoOpValidator, CharacterLoggedInHandler(l, span, worldId, channelId))
		})
	}
}

type characterLoggedInRequest struct {
	characterId uint32
}

func (c *characterLoggedInRequest) CharacterId() uint32 {
	return c.characterId
}

func readCharacterLoggedInRequest(reader *request.RequestReader) characterLoggedInRequest {
	cid := reader.ReadUint32()
	return characterLoggedInRequest{cid}
}

func CharacterLoggedInHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readCharacterLoggedInRequest(r)
		c, err := character.GetCharacterById(l, span)(p.CharacterId())
		if err != nil {
			return
		}

		s = session.SetAccountId(c.Attributes().AccountId())(s.SessionId())
		s = session.SetCharacterId(c.Attributes().Id())(s.SessionId())
		s = session.SetGm(c.Attributes().Gm())(s.SessionId())

		session.Login(l, span)(worldId, channelId, s.AccountId(), p.CharacterId())
		err = session.Announce(s, _map.WriteGetCharacterInfo(l)(channelId, c))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}

		keys, err := keymap.GetByCharacterId(l, span)(c.Attributes().Id())
		if err != nil || len(keys) == 0 {
			l.WithError(err).Warnf("Unable to send keybinding to character %d.", c.Attributes().Id())
		} else {
			err = session.Announce(s, keymap.WriteKeyMap(l)(keys))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		}
	}
}
