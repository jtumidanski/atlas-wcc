package shop

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	npcShopServicePrefix string = "/ms/nss/"
	npcShopService              = requests.BaseRequest + npcShopServicePrefix
	npcShopResource             = npcShopService + "npcs/%d/shop"
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

func requestShop(npcId uint32) Request {
	return makeRequest(fmt.Sprintf(npcShopResource, npcId))
}
