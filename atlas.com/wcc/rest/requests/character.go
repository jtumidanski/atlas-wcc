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
	characterEquippedItems             = charactersInventoryResource + "?type=equip&include=inventoryItems,equipmentStatistics"
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

func (c *character) GetEquippedItemsForCharacter(characterId uint32) (*attributes.InventoryDataContainer, error) {
	ar := &attributes.InventoryDataContainer{}
	err := get(fmt.Sprintf(characterEquippedItems, characterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}
