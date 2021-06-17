package portal

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	mapInformationServicePrefix string = "/ms/mis/"
	mapInformationService              = requests.BaseRequest + mapInformationServicePrefix
	mapsResource                       = mapInformationService + "maps/"
	portalsResource                    = mapsResource + "%d/portals"
	portalsByName                      = portalsResource + "?name=%s"
)

func requestPortalByName(mapId uint32, portalName string) (*PortalDataContainer, error) {
	ar := &PortalDataContainer{}
	err := requests.Get(fmt.Sprintf(portalsByName, mapId, portalName), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
