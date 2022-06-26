package wishlist

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	cashShopServicePrefix  string = "/ms/cashshop/"
	cashShopServiceService        = requests.BaseRequest + cashShopServicePrefix
	wishlistResource              = cashShopServiceService + "characters/%d/wishlist"
)

func requestById(characterId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(wishlistResource, characterId))
}
