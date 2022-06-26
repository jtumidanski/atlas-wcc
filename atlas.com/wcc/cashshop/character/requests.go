package character

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	cashShopServicePrefix  string = "/ms/cashshop/"
	cashShopServiceService        = requests.BaseRequest + cashShopServicePrefix
	charactersResource            = cashShopServiceService + "characters/"
	charactersById                = charactersResource + "%d"
)

func requestById(characterId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(charactersById, characterId))
}
