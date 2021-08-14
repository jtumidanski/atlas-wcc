package requests

import (
	"atlas-wcc/rest/attributes"
	"fmt"
	"github.com/sirupsen/logrus"
)

const (
	charactersServicePrefix     string = "/ms/cos/"
	charactersService                  = BaseRequest + charactersServicePrefix
	charactersResource                 = charactersService + "characters/"
	charactersById                     = charactersResource + "%d"
	charactersInventoryResource        = charactersResource + "%d/inventories/"
	characterItems                     = charactersInventoryResource + "?type=%s&include=inventoryItems,equipmentStatistics"
	characterItem                      = charactersInventoryResource + "?type=%s&slot=%d&include=inventoryItems,equipmentStatistics"
	characterWeaponDamage              = charactersResource + "%d/damage/weapon"
)

func GetCharacterAttributesById(l logrus.FieldLogger) func(characterId uint32) (*attributes.CharacterAttributesDataContainer, error) {
	return func(characterId uint32) (*attributes.CharacterAttributesDataContainer, error) {
		ar := &attributes.CharacterAttributesDataContainer{}
		err := Get(l)(fmt.Sprintf(charactersById, characterId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func GetItemsForCharacter(l logrus.FieldLogger) func(characterId uint32, inventoryType string) (*attributes.InventoryDataContainer, error) {
	return func(characterId uint32, inventoryType string) (*attributes.InventoryDataContainer, error) {
		ar := &attributes.InventoryDataContainer{}
		err := Get(l)(fmt.Sprintf(characterItems, characterId, inventoryType), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func GetItemForCharacter(l logrus.FieldLogger) func(characterId uint32, inventoryType string, slot int16) (*attributes.InventoryDataContainer, error) {
	return func(characterId uint32, inventoryType string, slot int16) (*attributes.InventoryDataContainer, error) {
		ar := &attributes.InventoryDataContainer{}
		err := Get(l)(fmt.Sprintf(characterItem, characterId, inventoryType, slot), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}

func GetEquippedItemsForCharacter(l logrus.FieldLogger) func(characterId uint32) (*attributes.InventoryDataContainer, error) {
	return func(characterId uint32) (*attributes.InventoryDataContainer, error) {
		return GetItemsForCharacter(l)(characterId, "equip")
	}
}

func GetEquippedItemForCharacter(l logrus.FieldLogger) func(characterId uint32, slot int16) (*attributes.InventoryDataContainer, error) {
	return func(characterId uint32, slot int16) (*attributes.InventoryDataContainer, error) {
		return GetItemForCharacter(l)(characterId, "equip", slot)
	}
}

func GetCharacterWeaponDamage(l logrus.FieldLogger) func(characterId uint32) (*attributes.DamageDataContainer, error) {
	return func(characterId uint32) (*attributes.DamageDataContainer, error) {
		ar := &attributes.DamageDataContainer{}
		err := Get(l)(fmt.Sprintf(characterWeaponDamage, characterId), ar)
		if err != nil {
			return nil, err
		}
		return ar, nil
	}
}
