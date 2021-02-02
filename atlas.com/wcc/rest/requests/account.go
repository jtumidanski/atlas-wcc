package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
)

const (
	accountsServicePrefix string = "/ms/aos/"
	accountsService              = baseRequest + accountsServicePrefix
	accountsResource             = accountsService + "accounts/"
	accountsById                 = accountsResource + "%d"
)

var Account = func() *account {
	return &account{}
}

type account struct {
}

func (a *account) GetById(id uint32) (*attributes.AccountDataContainer, error) {
	ar := &attributes.AccountDataContainer{}
	err := get(fmt.Sprintf(accountsById, id), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
