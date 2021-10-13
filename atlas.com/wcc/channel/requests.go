package channel

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	ServicePrefix string = "/ms/wrg/"
	Service              = requests.BaseRequest + ServicePrefix
	Resource             = Service + "channelServers/"
	ByWorld              = Resource + "?world=%d"
)

func requestForWorld(l logrus.FieldLogger, span opentracing.Span) func(worldId byte) (*dataContainer, error) {
	return func(worldId byte) (*dataContainer, error) {
		r := &dataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(ByWorld, worldId), r)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}
