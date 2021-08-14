package character

import (
	"atlas-wcc/character/properties"
	"atlas-wcc/character/skill"
	inventory2 "atlas-wcc/inventory"
	"atlas-wcc/pet"
)

type Model struct {
	attributes properties.Properties
	equipment  []inventory2.EquippedItem
	inventory  inventory2.Inventory
	skills     []skill.Model
	pets       []pet.Model
}

func (c Model) Attributes() properties.Properties {
	return c.attributes
}

func (c Model) Pets() []pet.Model {
	return c.pets
}

func (c Model) Equipment() []inventory2.EquippedItem {
	return c.equipment
}

func (c Model) Inventory() inventory2.Inventory {
	return c.inventory
}

func (c Model) Skills() []skill.Model {
	return c.skills
}

func (c Model) SetInventory(i inventory2.Inventory) Model {
	return Model{c.attributes, c.equipment, i, c.skills, c.pets}
}

func NewCharacter(attributes properties.Properties, equipment []inventory2.EquippedItem, skills []skill.Model, pets []pet.Model) Model {
	return Model{attributes, equipment, inventory2.EmptyInventory(), skills, pets}
}
