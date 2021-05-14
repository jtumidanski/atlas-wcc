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

func HasShop(l logrus.FieldLogger) func(npcId uint32) bool {
	return func(npcId uint32) bool {
		r, err := http.Get(fmt.Sprintf(npcShopResource, npcId))
		if err != nil {
			l.WithError(err).Errorf("Unable to identify if npc %d has a shop. Assuming not.", npcId)
			return false
		}
		return r.StatusCode == http.StatusOK
	}
}

func GetShop(l logrus.FieldLogger) func(npcId uint32) (*Model, error) {
	return func(npcId uint32) (*Model, error) {
		d := &dataContainer{}
		err := requests.Get(fmt.Sprintf(npcShopResource, npcId), d)
		if err != nil {
			return nil, err
		}

		return makeShop(d)
	}
}

func makeShop(d *dataContainer) (*Model, error) {
	items := make([]Item, 0)
	for _, i := range d.Data.Attributes.Items {
		items = append(items, Item{
			itemId:   i.ItemId,
			price:    i.Price,
			pitch:    i.Pitch,
			position: i.Position,
		})
	}

	return &Model{shopId: d.Data.Attributes.NPC, items: items}, nil
}
