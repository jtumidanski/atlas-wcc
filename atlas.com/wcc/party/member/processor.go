package member

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelOperator func(*Model)

type ModelListOperator func([]*Model)

type ModelProvider func() (*Model, error)

type ModelListProvider func() ([]*Model, error)

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

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r Request) ModelListProvider {
	return func(r Request) ModelListProvider {
		return func() ([]*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			ms := make([]*Model, 0)
			for _, v := range resp.DataList() {
				m, err := makeModel(&v)
				if err != nil {
					return nil, err
				}
				ms = append(ms, m)
			}
			return ms, nil
		}
	}
}

func InPartyModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(partyId uint32) ModelListProvider {
	return func(partyId uint32) ModelListProvider {
		return requestModelListProvider(l, span)(requestMembers(partyId))
	}
}

func GetInParty(l logrus.FieldLogger, span opentracing.Span) func(partyId uint32) ([]*Model, error) {
	return func(partyId uint32) ([]*Model, error) {
		return InPartyModelListProvider(l, span)(partyId)()
	}
}

func makeModel(body *dataBody) (*Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	att := body.Attributes
	m := Model{id: uint32(id), characterId: att.CharacterId, worldId: att.WorldId, channelId: att.ChannelId, online: att.Online}
	return &m, nil
}
