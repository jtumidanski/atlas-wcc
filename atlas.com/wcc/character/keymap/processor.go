package keymap

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ByCharacterIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) model.SliceProvider[Model] {
	return func(characterId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestKeyMap(characterId), makeModel)
	}
}

func GetByCharacterId(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		return ByCharacterIdModelProvider(l, span)(characterId)()
	}
}

func makeModel(k requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.Atoi(k.Id)
	if err != nil {
		return Model{}, err
	}

	attr := k.Attributes
	return Model{
		id:      uint32(id),
		key:     attr.Key,
		theType: attr.Type,
		action:  attr.Action,
	}, nil
}
