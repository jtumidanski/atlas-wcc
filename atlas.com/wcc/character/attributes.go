package character

import (
	"atlas-wcc/rest/requests"
	"atlas-wcc/rest/response"
)

type attributes struct {
	Type    string `json:"type"`
	Maximum uint32 `json:"maximum"`
}

const (
	ItemAttributesType      string = "com.atlas.cos.rest.attribute.ItemAttributes"
	EquipmentAttributesType string = "com.atlas.cos.rest.attribute.EquipmentAttributes"
	EquipmentStatisticsType string = "com.atlas.cos.rest.attribute.EquipmentStatisticsAttributes"
)

var equipmentIncludes = []response.ConditionalMapperProvider{
	transformItemAttributes,
	transformEquipmentAttributes,
	transformEquipmentStatistics,
}

func transformItemAttributes() (string, response.ObjectMapper) {
	return response.UnmarshalData(ItemAttributesType, func() interface{} {
		return requests.DataBody[itemAttributes]{}
	})
}

func transformEquipmentAttributes() (string, response.ObjectMapper) {
	return response.UnmarshalData(EquipmentAttributesType, func() interface{} {
		return requests.DataBody[equipmentAttributes]{}
	})
}

func transformEquipmentStatistics() (string, response.ObjectMapper) {
	return response.UnmarshalData(EquipmentStatisticsType, func() interface{} {
		return requests.DataBody[equipmentStatisticsAttributes]{}
	})
}

type inventoryAttributes struct {
	Type     string `json:"type"`
	Capacity byte   `json:"capacity"`
}

type itemAttributes struct {
	ItemId   uint32 `json:"itemId"`
	Quantity uint16 `json:"quantity"`
	Slot     int16  `json:"slot"`
}

type equipmentAttributes struct {
	EquipmentId int   `json:"equipmentId"`
	Slot        int16 `json:"slot"`
}

type equipmentStatisticsAttributes struct {
	ItemId        uint32 `json:"itemId"`
	Strength      uint16 `json:"strength"`
	Dexterity     uint16 `json:"dexterity"`
	Intelligence  uint16 `json:"intelligence"`
	Luck          uint16 `json:"luck"`
	Hp            uint16 `json:"hp"`
	Mp            uint16 `json:"mp"`
	WeaponAttack  uint16 `json:"weaponAttack"`
	MagicAttack   uint16 `json:"magicAttack"`
	WeaponDefense uint16 `json:"weaponDefense"`
	MagicDefense  uint16 `json:"magicDefense"`
	Accuracy      uint16 `json:"accuracy"`
	Avoidability  uint16 `json:"avoidability"`
	Hands         uint16 `json:"hands"`
	Speed         uint16 `json:"speed"`
	Jump          uint16 `json:"jump"`
	Slots         byte   `json:"slots"`
}
