package requests

import (
   "atlas-wcc/rest/attributes"
   "fmt"
)

const (
   MapInformationServicePrefix string = "/ms/mis/"
   MapInformationService              = BaseRequest + MapInformationServicePrefix
   MapsResource                       = MapInformationService + "maps/"
   NPCsInMap                          = MapsResource + "%d/npcs"
   PortalsResource                    = MapsResource + "%d/portals"
   PortalsByName                      = PortalsResource + "?name=%s"
)

func GetPortalByName(mapId uint32, portalName string) (*attributes.PortalDataContainer, error) {
   ar := &attributes.PortalDataContainer{}
   err := Get(fmt.Sprintf(PortalsByName, mapId, portalName), ar)
   if err != nil {
      return nil, err
   }
   return ar, nil
}

func GetNPCsInMap(mapId uint32) (*attributes.NpcDataContainer, error) {
   ar := &attributes.NpcDataContainer{}
   err := Get(fmt.Sprintf(NPCsInMap, mapId), ar)
   if err != nil {
      return nil, err
   }
   return ar, nil
}
