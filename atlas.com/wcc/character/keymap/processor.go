package keymap

import (
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelListProvider func() ([]*Model, error)

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r requests.Request[attributes]) ModelListProvider {
	return func(r requests.Request[attributes]) ModelListProvider {
		return func() ([]*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			ms := make([]*Model, 0)
			for _, v := range resp.DataList() {
				m, err := makeModel(v)
				if err != nil {
					return nil, err
				}
				ms = append(ms, m)
			}
			return ms, nil
		}
	}
}

func ByCharacterIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ModelListProvider {
	return func(characterId uint32) ModelListProvider {
		return requestModelListProvider(l, span)(requestKeyMap(characterId))
	}
}

func GetByCharacterId(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]*Model, error) {
	return func(characterId uint32) ([]*Model, error) {
		return ByCharacterIdModelProvider(l, span)(characterId)()
	}
}

func makeModel(k requests.DataBody[attributes]) (*Model, error) {
	id, err := strconv.Atoi(k.Id)
	if err != nil {
		return nil, err
	}

	attr := k.Attributes
	return &Model{
		id:      uint32(id),
		key:     attr.Key,
		theType: attr.Type,
		action:  attr.Action,
	}, nil
}
