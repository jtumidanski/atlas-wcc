package wishlist

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func byIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) model.SliceProvider[Model] {
	return func(characterId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestById(characterId), makeModel)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		return byIdModelProvider(l, span)(characterId)()
	}
}

func makeModel(db requests.DataBody[attributes]) (Model, error) {
	att := db.Attributes
	return Model{
		serialNumber: att.SerialNumber,
	}, nil
}
