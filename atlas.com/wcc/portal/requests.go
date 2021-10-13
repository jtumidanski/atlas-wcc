package portal

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	portalsResource                    = mapsResource + "%d/portals"
	portalsByName                      = portalsResource + "?name=%s"
)

func requestByName(l logrus.FieldLogger, span opentracing.Span) func(mapId uint32, portalName string) (*dataContainer, error) {
	return func(mapId uint32, portalName string) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(portalsByName, mapId, portalName), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
