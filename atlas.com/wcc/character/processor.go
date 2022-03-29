package character

import (
	"atlas-wcc/character/properties"
	"atlas-wcc/character/skill"
	"atlas-wcc/inventory"
	"atlas-wcc/pet"
	"atlas-wcc/rest/requests"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func GetCharacterById(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (*Model, error) {
	return func(characterId uint32) (*Model, error) {
		cs, err := properties.GetById(l, span)(characterId)
		if err != nil {
			return nil, err
		}

		c, err := getCharacterForAttributes(l, span)(cs)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func getCharacterForAttributes(l logrus.FieldLogger, span opentracing.Span) func(data *properties.Model) (*Model, error) {
	return func(data *properties.Model) (*Model, error) {
		eq, err := getEquippedItemsForCharacter(l, span)(data.Id())
		if err != nil {
			return nil, err
		}

		ps, err := getPetsForCharacter()
		if err != nil {
			return nil, err
		}

		ss, err := skill.GetForCharacter(l, span)(data.Id())
		if err != nil {
			return nil, err
		}

		c := NewCharacter(*data, eq, ss, ps)

		ei, err := getEquipInventoryForCharacter(l, span)(data.Id())
		if err != nil {
			return nil, err
		}
		ui, err := getItemInventoryForCharacter(l, span)(data.Id(), "use")
		if err != nil {
			return nil, err
		}
		si, err := getItemInventoryForCharacter(l, span)(data.Id(), "setup")
		if err != nil {
			return nil, err
		}
		etc, err := getItemInventoryForCharacter(l, span)(data.Id(), "etc")
		if err != nil {
			return nil, err
		}
		ci, err := getItemInventoryForCharacter(l, span)(data.Id(), "cash")
		if err != nil {
			return nil, err
		}
		i := c.Inventory().SetEquipInventory(*ei).SetUseInventory(*ui).SetSetupInventory(*si).SetEtcInventory(*etc).SetCashInventory(*ci)
		c = c.SetInventory(i)
		return &c, nil
	}
}

func getPetsForCharacter() ([]pet.Model, error) {
	return make([]pet.Model, 0), nil
}

func itemFilter(i requests.DataBody[itemAttributes]) bool {
	attr := i.Attributes
	return attr.Slot >= 0
}

func getItemInventoryForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, inventoryType string) (*inventory.ItemInventory, error) {
	return func(characterId uint32, inventoryType string) (*inventory.ItemInventory, error) {
		r, err := requestItemsForCharacter(characterId, inventoryType)(l, span)
		if err != nil {
			return nil, err
		}

		is := make([]inventory.Item, 0)
		for _, i := range requests.GetIncluded(r, itemFilter) {
			item := inventory.NewItem(i.ItemId, i.Slot, i.Quantity)
			is = append(is, item)
		}
		attr := r.Data().Attributes
		i := inventory.NewItemInventory(attr.Capacity, is)
		return &i, nil
	}
}

func getEquipInventoryForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (*inventory.EquipInventory, error) {
	return func(characterId uint32) (*inventory.EquipInventory, error) {
		r, err := requestItemsForCharacter(characterId, "equip")(l, span)
		if err != nil {
			return nil, err
		}

		eis := make([]inventory.EquippedItem, 0)
		for _, e := range requests.GetIncluded(r, equipmentItemFilter) {
			ea, ok := requests.GetInclude[inventoryAttributes, equipmentStatisticsAttributes](r, e.EquipmentId)
			if ok {
				ei := inventory.NewEquippedItemBuilder().
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
		ei := inventory.NewEquipInventory(attr.Capacity, eis)
		return &ei, nil
	}
}

func getEquippedItemsForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) ([]inventory.EquippedItem, error) {
	return func(characterId uint32) ([]inventory.EquippedItem, error) {
		r, err := requestEquippedItemsForCharacter(characterId)(l, span)
		if err != nil {
			return nil, err
		}

		eis := make([]inventory.EquippedItem, 0)
		for _, e := range requests.GetIncluded(r, equippedItemFilter) {
			ea, ok := requests.GetInclude[inventoryAttributes, equipmentStatisticsAttributes](r, e.EquipmentId)
			if ok {
				ei := inventory.NewEquippedItemBuilder().
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

func GetEquipItemForCharacter(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, slot int16) (*inventory.EquippedItem, error) {
	return func(characterId uint32, slot int16) (*inventory.EquippedItem, error) {
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
				ei := inventory.NewEquippedItemBuilder().
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

func GetCharacterWeaponDamage(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) uint32 {
	return func(characterId uint32) uint32 {
		r, err := requestCharacterWeaponDamage(characterId)(l, span)
		if err != nil {
			return 1
		}
		attr := r.Data().Attributes
		return attr.Maximum
	}
}
