package drop

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	dropRegistryServicePrefix string = "/ms/drg/"
	dropRegistryService              = requests.BaseRequest + dropRegistryServicePrefix
	dropResource                     = dropRegistryService + "worlds/%d/channels/%d/maps/%d/drops"
)

func requestInMap(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32) (*dataContainer, error) {
	return func(worldId byte, channelId byte, mapId uint32) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(dropResource, worldId, channelId, mapId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
