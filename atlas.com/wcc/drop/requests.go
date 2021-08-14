package drop

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	dropRegistryServicePrefix string = "/ms/drg/"
	dropRegistryService              = requests.BaseRequest + dropRegistryServicePrefix
	dropResource                     = dropRegistryService + "worlds/%d/channels/%d/maps/%d/drops"
)

func requestDropsInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) (*dataContainer, error) {
	return func(worldId byte, channelId byte, mapId uint32) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l)(fmt.Sprintf(dropResource, worldId, channelId, mapId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
