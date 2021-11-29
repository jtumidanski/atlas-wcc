package member

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	partyRegistryServicePrefix string = "/ms/party/"
	partyRegistryService              = requests.BaseRequest + partyRegistryServicePrefix
	partiesResource                   = partyRegistryService + "parties"
	partyResource                     = partiesResource + "/%d"
	partyMembersResource              = partyResource + "/members"
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

func requestMembers(partyId uint32) Request {
	return makeRequest(fmt.Sprintf(partyMembersResource, partyId))
}