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

var Account = func() *account {
	return &account{}
}

type account struct {
}

func (a *account) requestAccountById(id uint32) (*dataContainer, error) {
	ar := &dataContainer{}
	err := requests.Get(fmt.Sprintf(accountsById, id), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
