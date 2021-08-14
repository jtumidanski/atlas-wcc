package _map

import (
	"atlas-wcc/map/character"
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

func GetCharactersInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) (*character.MapCharacterDataContainer, error) {
	return func(worldId byte, channelId byte, mapId uint32) (*character.MapCharacterDataContainer, error) {
		ar := &character.MapCharacterDataContainer{}
		err := requests.Get(l)(fmt.Sprintf(mapCharactersResource, worldId, channelId, mapId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
