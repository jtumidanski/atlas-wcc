package domain

type Character struct {
   attributes CharacterAttributes
   equipment  []EquippedItem
   inventory  Inventory
   skills     []Skill
   pets       []Pet
}

func (c Character) Attributes() CharacterAttributes {
   return c.attributes
}

func (c Character) Pets() []Pet {
   return c.pets
}

func (c Character) Equipment() []EquippedItem {
   return c.equipment
}

func (c Character) Inventory() Inventory {
   return c.inventory
}

func (c Character) Skills() []Skill {
   return c.skills
}

func NewCharacter(attributes CharacterAttributes, equipment []EquippedItem, skills []Skill, pets []Pet) Character {
   return Character{attributes, equipment, EmptyInventory(), skills, pets}
}
