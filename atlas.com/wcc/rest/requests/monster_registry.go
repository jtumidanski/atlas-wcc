package requests

import (
   "atlas-wcc/rest/attributes"
   "fmt"
)

const (
   MonsterRegistryServicePrefix string = "/ms/morg/"
   MonsterRegistryService              = BaseRequest + MonsterRegistryServicePrefix
   MonsterResource                     = MonsterRegistryService + "worlds/%d/channels/%d/maps/%d/monsters"
)

func GetMonstersInMap(worldId byte, channelId byte, mapId uint32) (*attributes.MonsterDataContainer, error) {
   ar := &attributes.MonsterDataContainer{}
   err := Get(fmt.Sprintf(MonsterResource, worldId, channelId, mapId), ar)
   if err != nil {
      return nil, err
   }
   return ar, nil
}