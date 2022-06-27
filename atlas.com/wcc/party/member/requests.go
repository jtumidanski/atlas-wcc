package member

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	partyRegistryServicePrefix string = "/ms/party/"
	partyRegistryService              = requests.BaseRequest + partyRegistryServicePrefix
	partiesResource                   = partyRegistryService + "parties"
	partyResource                     = partiesResource + "/%d"
	partyMembersResource              = partyResource + "/members"
)

func requestMembers(partyId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(partyMembersResource, partyId))
}
