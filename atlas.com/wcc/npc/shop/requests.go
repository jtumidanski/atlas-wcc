package shop

import (
	"atlas-wcc/rest/requests"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
)

const (
	npcShopServicePrefix string = "/ms/nss/"
	npcShopService              = requests.BaseRequest + npcShopServicePrefix
	npcShopResource             = npcShopService + "npcs/%d/shop"
)

func hasShop(l logrus.FieldLogger) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		r, err := http.Get(fmt.Sprintf(npcShopResource, npcId))
		if err != nil {
			l.WithError(err).Errorf("Unable to identify if npc %d has a shop. Assuming not.", npcId)
			return false
		}
		return r.StatusCode == http.StatusOK
	}
}

func requestShop(l logrus.FieldLogger) func(npcId uint32) (*dataContainer, error) {
	return func(npcId uint32) (*dataContainer, error) {
		d := &dataContainer{}
		err := requests.Get(l)(fmt.Sprintf(npcShopResource, npcId), d)
		if err != nil {
			return nil, err
		}
		return d, nil
	}
}
