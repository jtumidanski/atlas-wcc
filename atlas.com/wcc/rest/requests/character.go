package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
)

const (
	CharactersServicePrefix     string = "/ms/cos/"
	CharactersService                  = BaseRequest + CharactersServicePrefix
	CharactersResource                 = CharactersService + "characters/"
	CharactersByName                   = CharactersResource + "?name=%s"
	CharactersForAccountByWorld        = CharactersResource + "?accountId=%d&worldId=%d"
	CharactersById                     = CharactersResource + "%d"
	CharactersInventoryResource        = CharactersResource + "%d/inventories/"
	CharacterEquippedItems             = CharactersInventoryResource + "?type=equip&include=inventoryItems,equipmentStatistics"
	CharacterSeeds                     = CharactersResource + "seeds/"
)

func GetCharacterAttributesByName(name string) (*attributes.CharacterAttributesDataContainer, error) {
	ar := &attributes.CharacterAttributesDataContainer{}
	err := Get(fmt.Sprintf(CharactersByName, name), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func GetCharacterAttributesForAccountByWorld(accountId uint32, worldId byte) (*attributes.CharacterAttributesDataContainer, error) {
	ar := &attributes.CharacterAttributesDataContainer{}
	err := Get(fmt.Sprintf(CharactersForAccountByWorld, accountId, worldId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func GetCharacterAttributesById(characterId uint32) (*attributes.CharacterAttributesDataContainer, error) {
	ar := &attributes.CharacterAttributesDataContainer{}
	err := Get(fmt.Sprintf(CharactersById, characterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}

func GetEquippedItemsForCharacter(characterId uint32) (*attributes.InventoryDataContainer, error) {
	ar := &attributes.InventoryDataContainer{}
	err := Get(fmt.Sprintf(CharacterEquippedItems, characterId), ar)
	if err != nil {
		return nil, err
	}
	return ar, nil
}