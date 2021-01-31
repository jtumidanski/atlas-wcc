package requests

import (
   "atlas-wcc/rest/attributes"
   "fmt"
)

const (
   MapRegistryServicePrefix string = "/ms/mrg/"
   MapRegistryService              = BaseRequest + MapRegistryServicePrefix
   MapResource                     = MapRegistryService + "worlds/%d/channels/%d/maps/%d"
   MapCharactersResource           = MapResource + "/characters/"
)

func GetCharactersInMap(worldId byte, channelId byte, mapId uint32) (*attributes.MapCharacterDataContainer, error) {
   ar := &attributes.MapCharacterDataContainer{}
   err := Get(fmt.Sprintf(MapCharactersResource, worldId, channelId, mapId), ar)
   if err != nil {
      return nil, err
   }
   return ar, nil
}
