package drop

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	dropRegistryServicePrefix string = "/ms/drg/"
	dropRegistryService              = requests.BaseRequest + dropRegistryServicePrefix
	dropResource                     = dropRegistryService + "worlds/%d/channels/%d/maps/%d/drops"
)

func requestInMap(worldId byte, channelId byte, mapId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(dropResource, worldId, channelId, mapId))
}
