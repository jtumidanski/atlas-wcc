package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
)

const (
	monsterRegistryServicePrefix string = "/ms/morg/"
	monsterRegistryService              = BaseRequest + monsterRegistryServicePrefix
	mapMonstersResource                 = monsterRegistryService + "worlds/%d/channels/%d/maps/%d/monsters"
	monstersResource                    = monsterRegistryService + "monsters"
	monsterResource                     = monstersResource + "/%d"
)

var MonsterRegistry = func() *monsterRegistry {
	return &monsterRegistry{}
}

type monsterRegistry struct {
}

func (m *monsterRegistry) GetInMap(worldId byte, channelId byte, mapId uint32) (*attributes.MonsterDataContainer, error) {
	ar := &attributes.MonsterDataContainer{}
	err := get(fmt.Sprintf(mapMonstersResource, worldId, channelId, mapId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (m *monsterRegistry) GetById(id uint32) (*attributes.MonsterDataContainer, error) {
	ar := &attributes.MonsterDataContainer{}
	err := get(fmt.Sprintf(monsterResource, id), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
