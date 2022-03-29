package channel

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	ServicePrefix string = "/ms/wrg/"
	Service              = requests.BaseRequest + ServicePrefix
	Resource             = Service + "channelServers/"
	ByWorld              = Resource + "?world=%d"
)

func requestForWorld(worldId byte) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(ByWorld, worldId))
}
