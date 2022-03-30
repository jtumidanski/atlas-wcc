package character

import (
	"atlas-wcc/character/inventory"
	"atlas-wcc/character/properties"
	"atlas-wcc/character/skill"
	"atlas-wcc/pet"
)

type Model struct {
	attributes properties.Model
	equipment  []inventory.EquippedItem
	inventory  inventory.Inventory
	skills     []skill.Model
	pets       []pet.Model
}

func (c Model) Attributes() properties.Model {
	return c.attributes
}

func (c Model) Pets() []pet.Model {
	return c.pets
}

func (c Model) Equipment() []inventory.EquippedItem {
	return c.equipment
}

func (c Model) Inventory() inventory.Inventory {
	return c.inventory
}

func (c Model) Skills() []skill.Model {
	return c.skills
}

func (c Model) SetInventory(i inventory.Inventory) Model {
	return Model{c.attributes, c.equipment, i, c.skills, c.pets}
}

func NewCharacter(attributes properties.Model, equipment []inventory.EquippedItem, skills []skill.Model, pets []pet.Model) Model {
	return Model{attributes, equipment, inventory.EmptyInventory(), skills, pets}
}
