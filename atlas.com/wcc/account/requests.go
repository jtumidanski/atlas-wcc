package account

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	accountsServicePrefix string = "/ms/aos/"
	accountsService              = requests.BaseRequest + accountsServicePrefix
	accountsResource             = accountsService + "accounts/"
	accountsById                 = accountsResource + "%d"
)

func requestAccountById(id uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(accountsById, id))
}
