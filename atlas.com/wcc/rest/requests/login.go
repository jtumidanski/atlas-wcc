package requests

import (
   "fmt"
)

const (
   LoginsResource = AccountsService + "logins/"
   LoginsById     = LoginsResource + "%d"
)

func CreateLogout(accountId uint32) {
   _, _ = Delete(fmt.Sprintf(LoginsById, accountId))
}
