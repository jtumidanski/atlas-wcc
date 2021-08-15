package reactor

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	reactorServicePrefix string = "/ms/ros/"
	reactorService              = requests.BaseRequest + reactorServicePrefix
	reactorsResource            = reactorService + "reactors/"
	reactorById                 = reactorsResource + "%d"
)

func requestById(l logrus.FieldLogger) func(id uint32) (*DataContainer, error) {
	return func(id uint32) (*DataContainer, error) {
		dc := &DataContainer{}
		err := requests.Get(l)(fmt.Sprintf(reactorById, id), dc)
		if err != nil {
			return nil, err
		}
		return dc, nil
	}
}
