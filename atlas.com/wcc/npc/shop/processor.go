package shop

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func HasShop(l logrus.FieldLogger) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		return hasShop(l)(npcId)
	}
}

func GetShop(l logrus.FieldLogger, span opentracing.Span) func(npcId uint32) (*Model, error) {
	return func(npcId uint32) (*Model, error) {
		d, err := requestShop(l, span)(npcId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve shop for %d.", npcId)
			return nil, err
		}
		return makeShop(d)
	}
}

func makeShop(d *dataContainer) (*Model, error) {
	items := make([]Item, 0)
	for _, i := range d.Data.Attributes.Items {
		items = append(items, Item{
			itemId:   i.ItemId,
			price:    i.Price,
			pitch:    i.Pitch,
			position: i.Position,
		})
	}

	return &Model{shopId: d.Data.Attributes.NPC, items: items}, nil
}