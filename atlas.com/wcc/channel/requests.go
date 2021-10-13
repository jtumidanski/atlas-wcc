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

type Request func(l logrus.FieldLogger, span opentracing.Span) (*dataContainer, error)

func makeRequest(url string) Request {
	return func(l logrus.FieldLogger, span opentracing.Span) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l, span)(url, ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func requestForWorld(worldId byte) Request {
	return makeRequest(fmt.Sprintf(ByWorld, worldId))
}
