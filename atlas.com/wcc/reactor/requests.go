package reactor

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	reactorServicePrefix string = "/ms/ros/"
	reactorService              = requests.BaseRequest + reactorServicePrefix
	reactorsResource            = reactorService + "reactors/"
	reactorById                 = reactorsResource + "%d"
	mapReactorsResource         = reactorService + "worlds/%d/channels/%d/maps/%d/reactors"
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

func requestInMap(worldId byte, channelId byte, mapId uint32) Request {
	return makeRequest(fmt.Sprintf(mapReactorsResource, worldId, channelId, mapId))
}

func requestById(id uint32) Request {
	return makeRequest(fmt.Sprintf(reactorById, id))
}
