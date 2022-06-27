package portal

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
)

type IdProvider func() uint32

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

func RandomPortalIdProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) IdProvider {
	return func(mapId uint32) IdProvider {
		return func() uint32 {
			ps, err := ForMap(l, span)(mapId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve portals for map %d. Defaulting to 0.", mapId)
				return 0
			}
			if len(ps) == 0 {
				l.Warnf("No portals in map %d. Defaulting to zero.", mapId)
				return 0
			}
			return ps[rand.Intn(len(ps))].Id()
		}
	}
}

func ForMap(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32) ([]Model, error) {
	return func(mapId uint32) ([]Model, error) {
		return requests.SliceProvider[attributes, Model](l, span)(requestAll(mapId), makeModel)()
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
