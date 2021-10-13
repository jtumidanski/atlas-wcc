package portal

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

type ModelProvider func() (*Model, error)

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

func ByNameModelProvider(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalName string) ModelProvider {
	return func(mapId uint32, portalName string) ModelProvider {
		return requestModelProvider(l, span)(requestByName(mapId, portalName))
	}
}

func GetByName(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalName string) (*Model, error) {
	return func(mapId uint32, portalName string) (*Model, error) {
		return ByNameModelProvider(l, span)(mapId, portalName)()
	}
}

func makeModel(body *dataBody) (*Model, error) {
	id, err := strconv.ParseUint(body.Id, 10, 32)
	if err != nil {
		return nil, err
	}
	attr := body.Attributes
	m := NewPortal(uint32(id), attr.Name, attr.Target, attr.TargetMapId, attr.Type, attr.X, attr.Y, attr.ScriptName)
	return &m, nil
}
