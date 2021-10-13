package account

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	accountsServicePrefix string = "/ms/aos/"
	accountsService              = requests.BaseRequest + accountsServicePrefix
	accountsResource             = accountsService + "accounts/"
	accountsById                 = accountsResource + "%d"
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

func requestAccountById(id uint32) Request {
	return makeRequest(fmt.Sprintf(accountsById, id))
}
