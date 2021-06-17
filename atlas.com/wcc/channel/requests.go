package channel

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	ServicePrefix   string = "/ms/wrg/"
	Service                = requests.BaseRequest + ServicePrefix
	ChannelResource        = Service + "channelServers/"
	ByWorld                = ChannelResource + "?world=%d"
)

func requestChannelsForWorld(worldId byte) (*dataContainer, error) {
	r := &dataContainer{}
	err := requests.Get(fmt.Sprintf(ByWorld, worldId), r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
