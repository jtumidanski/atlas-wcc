package keymap

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelListProvider func() ([]*Model, error)

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

func makeModel(k *dataBody) (*Model, error) {
	id, err := strconv.Atoi(k.Id)
	if err != nil {
		return nil, err
	}

	return &Model{
		id:      uint32(id),
		key:     k.Attributes.Key,
		theType: k.Attributes.Type,
		action:  k.Attributes.Action,
	}, nil
}
