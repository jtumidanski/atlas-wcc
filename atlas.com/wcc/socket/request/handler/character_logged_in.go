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

func CharacterLoggedInHandler(l logrus.FieldLogger, span opentracing.Span) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readCharacterLoggedInRequest(r)
		c, err := character.GetCharacterById(l, span)(p.CharacterId())
		if err != nil {
			return
		}

		s = session.SetAccountId(c.Attributes().AccountId())(s.SessionId())
		s = session.SetCharacterId(c.Attributes().Id())(s.SessionId())
		s = session.SetGm(c.Attributes().Gm())(s.SessionId())

		session.Login(l, span)(s.WorldId(), s.ChannelId(), s.AccountId(), p.CharacterId())
		err = session.Announce(_map.WriteGetCharacterInfo(l)(s.ChannelId(), c))(s)
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}

		keys, err := keymap.GetByCharacterId(l, span)(c.Attributes().Id())
		if err != nil || len(keys) == 0 {
			l.WithError(err).Warnf("Unable to send keybinding to character %d.", c.Attributes().Id())
		} else {
			err = session.Announce(keymap.WriteKeyMap(l)(keys))(s)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		}
	}
}
