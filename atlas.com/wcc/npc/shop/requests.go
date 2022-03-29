package shop

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	npcShopServicePrefix string = "/ms/nss/"
	npcShopService              = requests.BaseRequest + npcShopServicePrefix
	npcShopResource             = npcShopService + "npcs/%d/shop"
)

func requestShop(npcId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(npcShopResource, npcId))
}
