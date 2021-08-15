package reactor

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	reactorServicePrefix string = "/ms/ros/"
	reactorService              = requests.BaseRequest + reactorServicePrefix
	reactorsResource            = reactorService + "reactors/"
	reactorById                 = reactorsResource + "%d"
	mapReactorsResource         = reactorService + "worlds/%d/channels/%d/maps/%d/reactors"
)

func requestById(l logrus.FieldLogger) func(id uint32) (*DataContainer, error) {
	return func(id uint32) (*DataContainer, error) {
		dc := &DataContainer{}
		err := requests.Get(l)(fmt.Sprintf(reactorById, id), dc)
		if err != nil {
			return nil, err
		}
		return dc, nil
	}
}

func requestInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) (*DataListContainer, error) {
	return func(worldId byte, channelId byte, mapId uint32) (*DataListContainer, error) {
		dc := &DataListContainer{}
		err := requests.Get(l)(fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId), dc)
		if err != nil {
			return nil, err
		}
		return dc, nil
	}
}
