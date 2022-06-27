package inventory

import (
	"atlas-wcc/rest/requests"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func itemFilter(i requests.DataBody[itemAttributes]) bool {
	attr := i.Attributes
	return attr.Slot >= 0
}

func GetItemInventoryForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, inventoryType string) (*ItemInventory, error) {
	return func(characterId uint32, inventoryType string) (*ItemInventory, error) {
		r, err := requestItemsForCharacter(characterId, inventoryType)(l, span)
		if err != nil {
			return nil, err
		}

		is := make([]Item, 0)
		for _, i := range requests.GetIncluded(r, itemFilter) {
			item := NewItem(i.ItemId, i.Slot, i.Quantity)
			is = append(is, item)
		}
		attr := r.Data().Attributes
		i := NewItemInventory(attr.Capacity, is)
		return &i, nil
	}
}

func GetEquipInventoryForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (*EquipInventory, error) {
	return func(characterId uint32) (*EquipInventory, error) {
		r, err := requestItemsForCharacter(characterId, "equip")(l, span)
		if err != nil {
			return nil, err
		}

		eis := make([]EquippedItem, 0)
		for _, e := range requests.GetIncluded(r, equipmentItemFilter) {
			ea, ok := requests.GetInclude[inventoryAttributes, equipmentStatisticsAttributes](r, e.EquipmentId)
			if ok {
				ei := NewEquippedItemBuilder().
					SetItemId(ea.ItemId).
					SetSlot(e.Slot).
					SetStrength(ea.Strength).
					SetDexterity(ea.Dexterity).
					SetIntelligence(ea.Intelligence).
					SetLuck(ea.Luck).
					SetHp(ea.Hp).
					SetMp(ea.Mp).
					SetWeaponAttack(ea.WeaponAttack).
					SetMagicAttack(ea.MagicAttack).
					SetWeaponDefense(ea.WeaponDefense).
					SetMagicDefense(ea.MagicDefense).
					SetAccuracy(ea.Accuracy).
					SetAvoidability(ea.Avoidability).
					SetHands(ea.Hands).
					SetSpeed(ea.Speed).
					SetJump(ea.Jump).
					SetSlots(ea.Slots).
					Build()
				eis = append(eis, ei)
			}
		}
		attr := r.Data().Attributes
		ei := NewEquipInventory(attr.Capacity, eis)
		return &ei, nil
	}
}

func GetEquippedItemsForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]EquippedItem, error) {
	return func(characterId uint32) ([]EquippedItem, error) {
		r, err := requestEquippedItemsForCharacter(characterId)(l, span)
		if err != nil {
			return nil, err
		}

		eis := make([]EquippedItem, 0)
		for _, e := range requests.GetIncluded(r, equippedItemFilter) {
			ea, ok := requests.GetInclude[inventoryAttributes, equipmentStatisticsAttributes](r, e.EquipmentId)
			if ok {
				ei := NewEquippedItemBuilder().
					SetItemId(ea.ItemId).
					SetSlot(e.Slot).
					SetStrength(ea.Strength).
					SetDexterity(ea.Dexterity).
					SetIntelligence(ea.Intelligence).
					SetLuck(ea.Luck).
					SetHp(ea.Hp).
					SetMp(ea.Mp).
					SetWeaponAttack(ea.WeaponAttack).
					SetMagicAttack(ea.MagicAttack).
					SetWeaponDefense(ea.WeaponDefense).
					SetMagicDefense(ea.MagicDefense).
					SetAccuracy(ea.Accuracy).
					SetAvoidability(ea.Avoidability).
					SetHands(ea.Hands).
					SetSpeed(ea.Speed).
					SetJump(ea.Jump).
					SetSlots(ea.Slots).
					Build()
				eis = append(eis, ei)
			}
		}

		return eis, nil
	}
}

func equippedItemFilter(i requests.DataBody[equipmentAttributes]) bool {
	attr := i.Attributes
	return attr.Slot < 0
}

func equipmentItemFilter(i requests.DataBody[equipmentAttributes]) bool {
	attr := i.Attributes
	return attr.Slot >= 0
}

func GetEquipItemForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, slot int16) (*EquippedItem, error) {
	return func(characterId uint32, slot int16) (*EquippedItem, error) {
		r, err := requestEquippedItemForCharacter(characterId, slot)(l, span)
		if err != nil {
			return nil, err
		}

		var equips []equipmentAttributes
		if slot < 0 {
			equips = requests.GetIncluded(r, equippedItemFilter)
		} else {
			equips = requests.GetIncluded(r, equipmentItemFilter)
		}

		for _, e := range equips {
			ea, ok := requests.GetInclude[inventoryAttributes, equipmentStatisticsAttributes](r, e.EquipmentId)
			if ok {
				ei := NewEquippedItemBuilder().
					SetItemId(ea.ItemId).
					SetSlot(e.Slot).
					SetStrength(ea.Strength).
					SetDexterity(ea.Dexterity).
					SetIntelligence(ea.Intelligence).
					SetLuck(ea.Luck).
					SetHp(ea.Hp).
					SetMp(ea.Mp).
					SetWeaponAttack(ea.WeaponAttack).
					SetMagicAttack(ea.MagicAttack).
					SetWeaponDefense(ea.WeaponDefense).
					SetMagicDefense(ea.MagicDefense).
					SetAccuracy(ea.Accuracy).
					SetAvoidability(ea.Avoidability).
					SetHands(ea.Hands).
					SetSpeed(ea.Speed).
					SetJump(ea.Jump).
					SetSlots(ea.Slots).
					Build()
				return &ei, nil
			}
		}

		return nil, errors.New("equipment not found")
	}
}
