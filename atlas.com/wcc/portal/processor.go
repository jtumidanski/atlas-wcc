package portal

import (
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetPortalByName(l logrus.FieldLogger) func(mapId uint32, portalName string) (*Model, error) {
	return func(mapId uint32, portalName string) (*Model, error) {
		resp, err := requestPortalByName(l)(mapId, portalName)
		if err != nil {
			return nil, err
		}

		d := resp.Data()
		aid, err := strconv.ParseUint(d.Id, 10, 32)
		if err != nil {
			return nil, err
		}

		a := makePortal(uint32(aid), mapId, d.Attributes)
		return &a, nil
	}
}

func makePortal(id uint32, mapId uint32, attr attributes) Model {
	return NewPortal(id, mapId, attr.Name, attr.Target, attr.TargetMapId, attr.Type, attr.X, attr.Y, attr.ScriptName)
}
