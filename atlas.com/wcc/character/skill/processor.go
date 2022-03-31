package skill

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ByCharacterIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) model.SliceProvider[Model] {
	return func(characterId uint32) model.SliceProvider[Model] {
		return requests.SliceProvider[attributes, Model](l, span)(requestForCharacter(characterId), makeModel)
	}
}

func GetForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		return ByCharacterIdModelProvider(l, span)(characterId)()
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	att := body.Attributes
	return NewSkill(uint32(id), att.Level, att.MasterLevel, att.Expiration, false, false), nil
}
