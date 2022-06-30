package shop

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func HasShop(l logrus.FieldLogger, span opentracing.Span) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		_, err := ByNpcIdModelProvider(l, span)(npcId)()
		return err == nil
	}
}

func ByNpcIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(npcId uint32) model.Provider[Model] {
	return func(npcId uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestShop(npcId), makeModel)
	}
}

func GetByNpcId(l logrus.FieldLogger, span opentracing.Span) func(npcId uint32) (Model, error) {
	return func(npcId uint32) (Model, error) {
		return ByNpcIdModelProvider(l, span)(npcId)()
	}
}

func makeModel(d requests.DataBody[attributes]) (Model, error) {
	items := make([]Item, 0)
	attr := d.Attributes
	for _, i := range attr.Items {
		items = append(items, Item{
			itemId:   i.ItemId,
			price:    i.Price,
			pitch:    i.Pitch,
			position: i.Position,
		})
	}

	return Model{shopId: attr.NPC, items: items}, nil
}

func ShowShop(l logrus.FieldLogger, span opentracing.Span) func(npcId uint32) model.Operator[session.Model] {
	return func(npcId uint32) model.Operator[session.Model] {
		ns, err := GetByNpcId(l, span)(npcId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve shop for npc %d.", npcId)
			return model.ErrorOperator[session.Model](err)
		}
		return session.AnnounceOperator(WriteGetNPCShop(l)(ns))
	}
}
