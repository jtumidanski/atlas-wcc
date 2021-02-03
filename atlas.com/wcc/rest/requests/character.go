package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
)

const (
	charactersServicePrefix     string = "/ms/cos/"
	charactersService                  = baseRequest + charactersServicePrefix
	charactersResource                 = charactersService + "characters/"
	charactersById                     = charactersResource + "%d"
	charactersInventoryResource        = charactersResource + "%d/inventories/"
	characterItems                     = charactersInventoryResource + "?type=%s&include=inventoryItems,equipmentStatistics"
	characterWeaponDamage              = charactersResource + "%d/damage/weapon"
)

var Character = func() *character {
	return &character{}
}

type character struct {
}

func (c *character) GetCharacterAttributesById(characterId uint32) (*attributes.CharacterAttributesDataContainer, error) {
	ar := &attributes.CharacterAttributesDataContainer{}
	err := get(fmt.Sprintf(charactersById, characterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (c *character) GetItemsForCharacter(characterId uint32, inventoryType string) (*attributes.InventoryDataContainer, error) {
	ar := &attributes.InventoryDataContainer{}
	err := get(fmt.Sprintf(characterItems, characterId, inventoryType), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func (c *character) GetEquippedItemsForCharacter(characterId uint32) (*attributes.InventoryDataContainer, error) {
	return c.GetItemsForCharacter(characterId, "equip")
}

func (c *character) GetCharacterWeaponDamage(characterId uint32) (*attributes.DamageDataContainer, error) {
	ar := &attributes.DamageDataContainer{}
	err := get(fmt.Sprintf(characterWeaponDamage, characterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
