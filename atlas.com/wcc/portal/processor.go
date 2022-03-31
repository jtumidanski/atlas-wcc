package portal

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func ByNameModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalName string) model.Provider[Model] {
	return func(mapId uint32, portalName string) model.Provider[Model] {
		return requests.Provider[attributes, Model](l, span)(requestByName(mapId, portalName), makeModel)
	}
}

func GetByName(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalName string) (Model, error) {
	return func(mapId uint32, portalName string) (Model, error) {
		return ByNameModelProvider(l, span)(mapId, portalName)()
	}
}

func makeModel(body requests.DataBody[attributes]) (Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return Model{}, err
	}
	attr := body.Attributes
	m := NewPortal(uint32(id), attr.Name, attr.Target, attr.TargetMapId, attr.Type, attr.X, attr.Y, attr.ScriptName)
	return m, nil
}
