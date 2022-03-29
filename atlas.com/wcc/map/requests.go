package _map

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	mapRegistryServicePrefix string = "/ms/mrg/"
	mapRegistryService              = requests.BaseRequest + mapRegistryServicePrefix
	mapResource                     = mapRegistryService + "worlds/%d/channels/%d/maps/%d"
	mapCharactersResource           = mapResource + "/characters/"
)

func requestCharactersInMap(worldId byte, channelId byte, mapId uint32) requests.Request[characterAttributes] {
	return requests.MakeGetRequest[characterAttributes](fmt.Sprintf(mapCharactersResource, worldId, channelId, mapId))
}
