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

func requestItemsForCharacter(characterId uint32, inventoryType string) requests.Request[inventoryAttributes] {
	return requests.MakeGetRequest[inventoryAttributes](fmt.Sprintf(characterItems, characterId, inventoryType), requests.AddMappers(equipmentIncludes))
}

func requestItemForCharacter(characterId uint32, inventoryType string, slot int16) requests.Request[inventoryAttributes] {
	return requests.MakeGetRequest[inventoryAttributes](fmt.Sprintf(characterItem, characterId, inventoryType, slot), requests.AddMappers(equipmentIncludes))
}

func requestEquippedItemsForCharacter(characterId uint32) requests.Request[inventoryAttributes] {
	return requestItemsForCharacter(characterId, "equip")
}

func requestEquippedItemForCharacter(characterId uint32, slot int16) requests.Request[inventoryAttributes] {
	return requestItemForCharacter(characterId, "equip", slot)
}

func requestCharacterWeaponDamage(characterId uint32) requests.Request[attributes] {
	return requests.MakeGetRequest[attributes](fmt.Sprintf(characterWeaponDamage, characterId))
}
