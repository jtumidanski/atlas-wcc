package account

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	accountsServicePrefix string = "/ms/aos/"
	accountsService              = requests.BaseRequest + accountsServicePrefix
	accountsResource             = accountsService + "accounts/"
	accountsById                 = accountsResource + "%d"
)

func requestAccountById(l logrus.FieldLogger) func(id uint32) (*dataContainer, error) {
	return func(id uint32) (*dataContainer, error) {
		ar := &dataContainer{}
		err := requests.Get(l)(fmt.Sprintf(accountsById, id), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
