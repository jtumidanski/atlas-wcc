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

func requestByName(mapId uint32, portalName string) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(portalsByName, mapId, portalName))
}
