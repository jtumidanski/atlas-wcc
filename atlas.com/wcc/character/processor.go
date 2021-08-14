package character

import (
	"atlas-wcc/character/properties"
	"atlas-wcc/character/skill"
	"atlas-wcc/inventory"
	"atlas-wcc/map"
	"atlas-wcc/pet"
	"atlas-wcc/rest/attributes"
	"atlas-wcc/rest/requests"
	"errors"
	"github.com/sirupsen/logrus"
	"strconv"
)

func GetCharacterById(l logrus.FieldLogger) func(characterId uint32) (*Model, error) {
	return func(characterId uint32) (*Model, error) {
		cs, err := properties.GetById(l)(characterId)
		if err != nil {
			return nil, err
		}

		c, err := getCharacterForAttributes(l)(cs)
		if err != nil {
			return nil, err
		}
		return c, nil
	}
}

func getCharacterForAttributes(l logrus.FieldLogger) func(data *properties.Properties) (*Model, error) {
	return func(data *properties.Properties) (*Model, error) {
		eq, err := getEquippedItemsForCharacter(l)(data.Id())
		if err != nil {
			return nil, err
		}

		ps, err := getPetsForCharacter()
		if err != nil {
			return nil, err
		}

		ss, err := skill.GetForCharacter(l)(data.Id())
		if err != nil {
			return nil, err
		}

		c := NewCharacter(*data, eq, ss, ps)

		ei, err := getEquipInventoryForCharacter(l)(data.Id())
		if err != nil {
			return nil, err
		}
		ui, err := getItemInventoryForCharacter(l)(data.Id(), "use")
		if err != nil {
			return nil, err
		}
		si, err := getItemInventoryForCharacter(l)(data.Id(), "setup")
		if err != nil {
			return nil, err
		}
		etc, err := getItemInventoryForCharacter(l)(data.Id(), "etc")
		if err != nil {
			return nil, err
		}
		ci, err := getItemInventoryForCharacter(l)(data.Id(), "cash")
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

func getItemInventoryForCharacter(l logrus.FieldLogger) func(characterId uint32, inventoryType string) (*inventory.ItemInventory, error) {
	return func(characterId uint32, inventoryType string) (*inventory.ItemInventory, error) {
		r, err := requests.GetItemsForCharacter(l)(characterId, inventoryType)
		if err != nil {
			return nil, err
		}

		is := make([]inventory.Item, 0)
		for _, i := range r.GetIncludedItems() {
			item := inventory.NewItem(i.Attributes.ItemId, i.Attributes.Slot, i.Attributes.Quantity)
			is = append(is, item)
		}
		i := inventory.NewItemInventory(r.Data().Attributes.Capacity, is)
		return &i, nil
	}
}

func getEquipInventoryForCharacter(l logrus.FieldLogger) func(characterId uint32) (*inventory.EquipInventory, error) {
	return func(characterId uint32) (*inventory.EquipInventory, error) {
		r, err := requests.GetItemsForCharacter(l)(characterId, "equip")
		if err != nil {
			return nil, err
		}

		eis := make([]inventory.EquippedItem, 0)
		for _, e := range r.GetIncludedEquips() {
			ea := r.GetEquipmentStatistics(e.Attributes.EquipmentId)
			if ea != nil {
				ei := inventory.NewEquippedItemBuilder().
					SetItemId(ea.ItemId).
					SetSlot(e.Attributes.Slot).
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

		ei := inventory.NewEquipInventory(r.Data().Attributes.Capacity, eis)
		return &ei, nil
	}
}

func getEquippedItemsForCharacter(l logrus.FieldLogger) func(characterId uint32) ([]inventory.EquippedItem, error) {
	return func(characterId uint32) ([]inventory.EquippedItem, error) {
		r, err := requests.GetEquippedItemsForCharacter(l)(characterId)
		if err != nil {
			return nil, err
		}

		eis := make([]inventory.EquippedItem, 0)
		for _, e := range r.GetIncludedEquippedItems() {
			ea := r.GetEquipmentStatistics(e.Attributes.EquipmentId)
			if ea != nil {
				ei := inventory.NewEquippedItemBuilder().
					SetItemId(ea.ItemId).
					SetSlot(e.Attributes.Slot).
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

func GetEquipItemForCharacter(l logrus.FieldLogger) func(characterId uint32, slot int16) (*inventory.EquippedItem, error) {
	return func(characterId uint32, slot int16) (*inventory.EquippedItem, error) {
		r, err := requests.GetEquippedItemForCharacter(l)(characterId, slot)
		if err != nil {
			return nil, err
		}

		var equips []attributes.EquipmentData
		if slot < 0 {
			equips = r.GetIncludedEquippedItems()
		} else {
			equips = r.GetIncludedEquips()
		}

		for _, e := range equips {
			ea := r.GetEquipmentStatistics(e.Attributes.EquipmentId)
			if ea != nil {
				ei := inventory.NewEquippedItemBuilder().
					SetItemId(ea.ItemId).
					SetSlot(e.Attributes.Slot).
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

func GetCharacterWeaponDamage(l logrus.FieldLogger) func(characterId uint32) uint32 {
	return func(characterId uint32) uint32 {
		r, err := requests.GetCharacterWeaponDamage(l)(characterId)
		if err != nil {
			return 1
		}
		return r.Data().Attributes.Maximum
	}
}

func GetCharacterIdsInMap(l logrus.FieldLogger) func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
	return func(worldId byte, channelId byte, mapId uint32) ([]uint32, error) {
		resp, err := _map.GetCharactersInMap(l)(worldId, channelId, mapId)
		if err != nil {
			return nil, err
		}

		cIds := make([]uint32, 0)
		for _, d := range resp.DataList() {
			cId, err := strconv.ParseUint(d.Id, 10, 32)
			if err != nil {
				break
			}
			cIds = append(cIds, uint32(cId))
		}
		return cIds, nil
	}
}
