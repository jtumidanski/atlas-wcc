package requests

import (
   "atlas-wcc/rest/attributes"
   "fmt"
)

const (
   MapInformationServicePrefix string = "/ms/mis/"
   MapInformationService              = BaseRequest + MapInformationServicePrefix
   MapsResource                       = MapInformationService + "maps/"
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
