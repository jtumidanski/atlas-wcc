package character

import (
	"atlas-wcc/rest/requests"
	"fmt"
)

const (
	charactersServicePrefix     string = "/ms/cos/"
	charactersService                  = requests.BaseRequest + charactersServicePrefix
	charactersResource                 = charactersService + "characters/"
	charactersById                     = charactersResource + "%d"
	charactersInventoryResource        = charactersResource + "%d/inventories/"
	characterItems                     = charactersInventoryResource + "?type=%s&include=inventoryItems,equipmentStatistics"
	characterItem                      = charactersInventoryResource + "?type=%s&slot=%d&include=inventoryItems,equipmentStatistics"
	characterWeaponDamage              = charactersResource + "%d/damage/weapon"
)

func requestCharacterWeaponDamage(characterId uint32) requests.Request[damageAttributes] {
	return requests.MakeGetRequest[damageAttributes](fmt.Sprintf(characterWeaponDamage, characterId))
}
