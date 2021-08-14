package channel

import (
	"errors"
	"github.com/sirupsen/logrus"
)

func GetForWorld(l logrus.FieldLogger) func(worldId byte, channelId byte) (*Model, error) {
	return func(worldId byte, channelId byte) (*Model, error) {
		r, err := requestChannelsForWorld(l)(worldId)
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
}

func makeChannel(data dataBody) Model {
	att := data.Attributes
	return NewChannelBuilder().
		SetWorldId(att.WorldId).
		SetChannelId(att.ChannelId).
		SetCapacity(att.Capacity).
		SetIpAddress(att.IpAddress).
		SetPort(att.Port).
		Build()
}
