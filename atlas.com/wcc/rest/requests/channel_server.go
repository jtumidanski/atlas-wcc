package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
)

const (
	ChannelServersResource            = worldRegistryService + "channelServers/"
	ChannelServersByWorld             = ChannelServersResource + "?world=%d"
)

func GetChannelsForWorld(worldId byte) (*attributes.ChannelServerDataContainer, error) {
	r := &attributes.ChannelServerDataContainer{}
	err := Get(fmt.Sprintf(ChannelServersByWorld, worldId), r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
