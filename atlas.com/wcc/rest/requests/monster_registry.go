package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
)

const (
	MonsterRegistryServicePrefix string = "/ms/morg/"
	MonsterRegistryService              = BaseRequest + MonsterRegistryServicePrefix
	MapMonstersResource                 = MonsterRegistryService + "worlds/%d/channels/%d/maps/%d/monsters"
	MonstersResource                    = MonsterRegistryService + "monsters"
	MonsterResource                     = MonstersResource + "/%d"
)

func GetMonstersInMap(worldId byte, channelId byte, mapId uint32) (*attributes.MonsterDataContainer, error) {
	ar := &attributes.MonsterDataContainer{}
	err := Get(fmt.Sprintf(MapMonstersResource, worldId, channelId, mapId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func GetMonster(id uint32) (*attributes.MonsterDataContainer, error) {
	ar := &attributes.MonsterDataContainer{}
	err := Get(fmt.Sprintf(MonsterResource, id), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
