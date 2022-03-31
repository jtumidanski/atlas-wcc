package channel

import (
	"atlas-wcc/model"
	"atlas-wcc/rest/requests"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func idFilter(channelId byte) func(models []Model) (Model, error) {
	return func(models []Model) (Model, error) {
		for _, m := range models {
			if m.ChannelId() == channelId {
				return m, nil
			}
		}
		return Model{}, errors.New("unable to locate channel for world")
	}
}

func ByWorldIdModelProvider(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte) model.Provider[Model] {
	return func(worldId byte, channelId byte) model.Provider[Model] {
		sp := requests.SliceProvider[attributes, Model](l, span)(requestForWorld(worldId), makeModel)
		return model.ModelListProviderToModelProviderAdapter(sp, idFilter(channelId))
	}
}

func GetByWorldId(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte) (Model, error) {
	return func(worldId byte, channelId byte) (Model, error) {
		return ByWorldIdModelProvider(l, span)(worldId, channelId)()
	}
}

func makeModel(data requests.DataBody[attributes]) (Model, error) {
	att := data.Attributes
	c := NewBuilder().
		SetWorldId(att.WorldId).
		SetChannelId(att.ChannelId).
		SetCapacity(att.Capacity).
		SetIpAddress(att.IpAddress).
		SetPort(att.Port).
		Build()
	return c, nil
}
