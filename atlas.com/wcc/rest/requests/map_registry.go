package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
)

const (
	mapRegistryServicePrefix string = "/ms/mrg/"
	mapRegistryService              = BaseRequest + mapRegistryServicePrefix
	mapResource                     = mapRegistryService + "worlds/%d/channels/%d/maps/%d"
	mapCharactersResource           = mapResource + "/characters/"
)

var MapRegistry = func() *mapRegistry {
	return &mapRegistry{}
}

type mapRegistry struct {
}

func (m *mapRegistry) GetCharactersInMap(worldId byte, channelId byte, mapId uint32) (*attributes.MapCharacterDataContainer, error) {
	ar := &attributes.MapCharacterDataContainer{}
	err := Get(fmt.Sprintf(mapCharactersResource, worldId, channelId, mapId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
