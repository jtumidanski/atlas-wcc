package character

import (
	"atlas-wcc/character/inventory"
	"atlas-wcc/character/properties"
	"atlas-wcc/character/skill"
	"atlas-wcc/pet"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

func GetCharacterById(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32) (Model, error) {
	return func(characterId uint32) (Model, error) {
		cs, err := properties.GetById(l, span)(characterId)
		if err != nil {
			return Model{}, err
		}

		c, err := getCharacterForAttributes(l, span)(cs)
		if err != nil {
			return Model{}, err
		}
		return c, nil
	}
}

func getCharacterForAttributes(l logrus.FieldLogger, span opentracing.Span) func(data properties.Model) (Model, error) {
	return func(data properties.Model) (Model, error) {
		eq, err := inventory.GetEquippedItemsForCharacter(l, span)(data.Id())
		if err != nil {
			return Model{}, err
		}

		ps, err := getPetsForCharacter()
		if err != nil {
			return Model{}, err
		}

		ss, err := skill.GetForCharacter(l, span)(data.Id())
		if err != nil {
			return Model{}, err
		}

		c := NewCharacter(data, eq, ss, ps)

		ei, err := inventory.GetEquipInventoryForCharacter(l, span)(data.Id())
		if err != nil {
			return Model{}, err
		}
		ui, err := inventory.GetItemInventoryForCharacter(l, span)(data.Id(), "use")
		if err != nil {
			return Model{}, err
		}
		si, err := inventory.GetItemInventoryForCharacter(l, span)(data.Id(), "setup")
		if err != nil {
			return Model{}, err
		}
		etc, err := inventory.GetItemInventoryForCharacter(l, span)(data.Id(), "etc")
		if err != nil {
			return Model{}, err
		}
		ci, err := inventory.GetItemInventoryForCharacter(l, span)(data.Id(), "cash")
		if err != nil {
			return Model{}, err
		}
		i := c.Inventory().SetEquipInventory(*ei).SetUseInventory(*ui).SetSetupInventory(*si).SetEtcInventory(*etc).SetCashInventory(*ci)
		c = c.SetInventory(i)
		return c, nil
	}
}

func getPetsForCharacter() ([]pet.Model, error) {
	return make([]pet.Model, 0), nil
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

func GainMeso(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, amount int32) {
	adjuster := emitMesoAdjustment(l, span)
	return func(characterId uint32, amount int32) {
		adjuster(characterId, amount)
	}
}
