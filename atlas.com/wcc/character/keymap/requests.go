package keymap

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	keymapServicePrefix string = "/ms/cks/"
	keymapService              = requests.BaseRequest + keymapServicePrefix
	keymapResource             = keymapService + "characters/%d/keymap"
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

func requestKeyMap(characterId uint32) Request {
	return makeRequest(fmt.Sprintf(keymapResource, characterId))
}
