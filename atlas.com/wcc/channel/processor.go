package channel

import (
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type ModelProvider func() (*Model, error)

type ModelListProvider func() ([]*Model, error)

func requestModelListProvider(l logrus.FieldLogger, span opentracing.Span) func(r Request) ModelListProvider {
	return func(r Request) ModelListProvider {
		return func() ([]*Model, error) {
			resp, err := r(l, span)
			if err != nil {
				return nil, err
			}

			ms := make([]*Model, 0)
			for _, v := range resp.DataList() {
				m, err := makeModel(&v)
				if err != nil {
					return nil, err
				}
				ms = append(ms, m)
			}
			return ms, nil
		}
	}
}

func ByWorldIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte) ModelProvider {
	return func(worldId byte, channelId byte) ModelProvider {
		return func() (*Model, error) {
			ms, err := requestModelListProvider(l, span)(requestForWorld(worldId))()
			if err != nil {
				return nil, err
			}
			for _, m := range ms {
				if m.ChannelId() == channelId {
					return m, nil
				}
			}
			return nil, errors.New("unable to locate channel for world")
		}
	}
}

func GetByWorldId(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte) (*Model, error) {
	return func(worldId byte, channelId byte) (*Model, error) {
		return ByWorldIdModelProvider(l, span)(worldId, channelId)()
	}
}

func makeModel(data *dataBody) (*Model, error) {
	att := data.Attributes
	c := NewBuilder().
		SetWorldId(att.WorldId).
		SetChannelId(att.ChannelId).
		SetCapacity(att.Capacity).
		SetIpAddress(att.IpAddress).
		SetPort(att.Port).
		Build()
	return &c, nil
}
