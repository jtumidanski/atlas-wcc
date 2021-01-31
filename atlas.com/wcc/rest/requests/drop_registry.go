package requests

import (
   "atlas-wcc/rest/attributes"
   "fmt"
)

const (
   DropRegistryServicePrefix string = "/ms/drg/"
   DropRegistryService              = BaseRequest + DropRegistryServicePrefix
   DropResource                     = DropRegistryService + "worlds/%d/channels/%d/maps/%d/drops"
)

func GetDropsInMap(worldId byte, channelId byte, mapId uint32) (*attributes.DropDataContainer, error) {
   ar := &attributes.DropDataContainer{}
   err := Get(fmt.Sprintf(DropResource, worldId, channelId, mapId), ar)
   if err != nil {
      return nil, err
   }
   return ar, nil
}