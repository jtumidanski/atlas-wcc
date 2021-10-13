package shop

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

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

func HasShop(l logrus.FieldLogger, span opentracing.Span) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		m, err := ByNpcIdModelProvider(l, span)(npcId)()
		if err != nil {
			return false
		}
		return m != nil
	}
}

func ByNpcIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(npcId uint32) ModelProvider {
	return func(npcId uint32) ModelProvider {
		return requestModelProvider(l, span)(requestShop(npcId))
	}
}

func GetByNpcId(l logrus.FieldLogger, span opentracing.Span) func(npcId uint32) (*Model, error) {
	return func(npcId uint32) (*Model, error) {
		return ByNpcIdModelProvider(l, span)(npcId)()
	}
}

func makeModel(d *dataBody) (*Model, error) {
	items := make([]Item, 0)
	for _, i := range d.Attributes.Items {
		items = append(items, Item{
			itemId:   i.ItemId,
			price:    i.Price,
			pitch:    i.Pitch,
			position: i.Position,
		})
	}

	return &Model{shopId: d.Attributes.NPC, items: items}, nil
}
