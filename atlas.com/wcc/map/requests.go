package _map

import (
	"atlas-wcc/map/character"
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	mapRegistryServicePrefix string = "/ms/mrg/"
	mapRegistryService              = requests.BaseRequest + mapRegistryServicePrefix
	mapResource                     = mapRegistryService + "worlds/%d/channels/%d/maps/%d"
	mapCharactersResource           = mapResource + "/characters/"
)

var MapRegistry = func() *mapRegistry {
	return &mapRegistry{}
}

type mapRegistry struct {
}

func (m *mapRegistry) GetCharactersInMap(worldId byte, channelId byte, mapId uint32) (*character.MapCharacterDataContainer, error) {
	ar := &character.MapCharacterDataContainer{}
	err := requests.Get(fmt.Sprintf(mapCharactersResource, worldId, channelId, mapId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}