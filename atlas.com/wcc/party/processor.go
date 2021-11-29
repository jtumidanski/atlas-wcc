package party

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelProvider func() (*Model, error)

func requestModelProvider(l logrus.FieldLogger, span opentracing.Span) func(r Request) ModelProvider {
	return func(r Request) ModelProvider {
		return func() (*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			p, err := makeModel(resp.Data())
			if err != nil {
				return nil, err
			}
			return p, nil
		}
	}
}

func ByMemberIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(memberId uint32) ModelProvider {
	return func(memberId uint32) ModelProvider {
		return requestModelProvider(l, span)(requestByMemberId(memberId))
	}
}

func GetByMemberId(l logrus.FieldLogger, span opentracing.Span) func(memberId uint32) (*Model, error) {
	return func(memberId uint32) (*Model, error) {
		return ByMemberIdModelProvider(l, span)(memberId)()
	}
}

func Create(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		err := createParty(l, span)(worldId, channelId, characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to create party for character %d.", characterId)
			return
		}
	}
}

func Leave(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, characterId uint32) {
	return func(worldId byte, channelId byte, characterId uint32) {
		p, err := GetByMemberId(l, span)(characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to locate party for character %d.", characterId)
			return
		}

		err = leaveParty(l, span)(worldId, channelId, p.Id(), characterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to create party for character %d.", characterId)
			return
		}
	}
}

func makeModel(body *dataBody) (*Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	att := body.Attributes
	m := Model{id: uint32(id), leaderId: att.LeaderId}
	return &m, nil
}
