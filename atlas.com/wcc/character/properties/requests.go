package properties

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	charactersServicePrefix     string = "/ms/cos/"
	charactersService                  = requests.BaseRequest + charactersServicePrefix
	charactersResource                 = charactersService + "characters/"
	charactersById                     = charactersResource + "%d"
)

func requestPropertiesById(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (*DataContainer, error) {
	return func(characterId uint32) (*DataContainer, error) {
		ar := &DataContainer{}
		err := requests.Get(l, span)(fmt.Sprintf(charactersById, characterId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}