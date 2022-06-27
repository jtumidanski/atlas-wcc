package reactor

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	reactorServicePrefix string = "/ms/ros/"
	reactorService              = requests.BaseRequest + reactorServicePrefix
	reactorsResource            = reactorService + "reactors/"
	reactorById                 = reactorsResource + "%d"
	mapReactorsResource         = reactorService + "worlds/%d/channels/%d/maps/%d/reactors"
)

func requestInMap(worldId byte, channelId byte, mapId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId))
}

func requestById(id uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(reactorById, id))
}
