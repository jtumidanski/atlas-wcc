package _map

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	mapRegistryServicePrefix string = "/ms/mrg/"
	mapRegistryService              = requests.BaseRequest + mapRegistryServicePrefix
	mapResource                     = mapRegistryService + "worlds/%d/channels/%d/maps/%d"
	mapCharactersResource           = mapResource + "/characters/"
)

func requestCharactersInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) (*CharacterDataContainer, error) {
	return func(worldId byte, channelId byte, mapId uint32) (*CharacterDataContainer, error) {
		ar := &CharacterDataContainer{}
		err := requests.Get(l)(fmt.Sprintf(mapCharactersResource, worldId, channelId, mapId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
