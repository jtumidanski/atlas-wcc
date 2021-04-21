package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	request2 "atlas-wcc/socket/request"
	"atlas-wcc/socket/response/writer"
	"context"
	"github.com/jtumidanski/atlas-socket/request"
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

func CharacterLoggedInHandler() request2.SessionRequestHandler {
	return func(l logrus.FieldLogger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readCharacterLoggedInRequest(r)
		c, err := processors.GetCharacterById(p.CharacterId())
		if err != nil {
			return
		}

		(*s).SetAccountId(c.Attributes().AccountId())
		(*s).SetCharacterId(c.Attributes().Id())
		(*s).SetGm(c.Attributes().Gm())

		producers.CharacterStatus(l, context.Background()).Login((*s).WorldId(), (*s).ChannelId(), (*s).AccountId(), p.CharacterId())
		(*s).Announce(writer.WriteGetCharacterInfo((*s).ChannelId(), *c))
	}
}
