package requests

import (
   "atlas-wcc/rest/attributes"
   "fmt"
)

const (
   dropRegistryServicePrefix string = "/ms/drg/"
   dropRegistryService = BaseRequest + dropRegistryServicePrefix
   dropResource        = dropRegistryService + "worlds/%d/channels/%d/maps/%d/drops"
)

var DropRegistry = func() *dropRegistry {
   return &dropRegistry{}
}

type dropRegistry struct {
}

func (d *dropRegistry) GetDropsInMap(worldId byte, channelId byte, mapId uint32) (*attributes.DropDataContainer, error) {
   ar := &attributes.DropDataContainer{}
   err := get(fmt.Sprintf(dropResource, worldId, channelId, mapId), ar)
   if err != nil {
      return nil, err
   }
   return ar, nil
}