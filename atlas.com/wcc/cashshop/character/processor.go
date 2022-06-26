package character

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func byIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) model.Provider[Model] {
	return func(characterId uint32) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestById(characterId), makeModel)
	}
}

func GetById(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		return byIdModelProvider(l, span)(characterId)()
	}
}

func makeModel(db requests.DataBody[attributes]) (Model, error) {
	cid, err := strconv.ParseUint(db.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	att := db.Attributes
	return Model{
		characterId: uint32(cid),
		credit:      att.Credit,
		points:      att.Points,
		prepaid:     att.Prepaid,
	}, nil
}
