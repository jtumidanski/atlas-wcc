package processors

import (
	"atlas-wcc/domain"
	"atlas-wcc/rest/attributes"
	"atlas-wcc/rest/requests"
	"errors"
)

func GetChannelForWorld(worldId byte, channelId byte) (*domain.Channel, error) {
	r, err := requests.GetChannelsForWorld(worldId)
	if err != nil {
		return nil, err
	}

	for _, x := range r.DataList() {
		w := makeChannel(x)
		if w.ChannelId() == channelId {
			return &w, nil
		}
	}
	return nil, errors.New("unable to locate channel for world")
}

func makeChannel(data attributes.ChannelServerData) domain.Channel {
	att := data.Attributes
	return domain.NewChannelBuilder().
		SetWorldId(att.WorldId).
		SetChannelId(att.ChannelId).
		SetCapacity(att.Capacity).
		SetIpAddress(att.IpAddress).
		SetPort(att.Port).
		Build()
}
