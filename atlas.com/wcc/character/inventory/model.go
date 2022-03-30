package inventory

type Inventory struct {
	equipInventory EquipInventory
	useInventory   ItemInventory
	setupInventory ItemInventory
	etcInventory   ItemInventory
	cashInventory  ItemInventory
}

func EmptyInventory() Inventory {
	return Inventory{
		equipInventory: EmptyEquipInventory(),
		useInventory:   EmptyItemInventory(),
		setupInventory: EmptyItemInventory(),
		etcInventory:   EmptyItemInventory(),
		cashInventory:  EmptyItemInventory(),
	}
}

func (i Inventory) SetEquipInventory(ei EquipInventory) Inventory {
	return Inventory{
		equipInventory: ei,
		useInventory:   i.useInventory,
		setupInventory: i.setupInventory,
		etcInventory:   i.etcInventory,
		cashInventory:  i.cashInventory,
	}
}

func (i Inventory) SetUseInventory(ii ItemInventory) Inventory {
	return Inventory{
		equipInventory: i.equipInventory,
		useInventory:   ii,
		setupInventory: i.setupInventory,
		etcInventory:   i.etcInventory,
		cashInventory:  i.cashInventory,
	}
}

func (i Inventory) SetSetupInventory(ii ItemInventory) Inventory {
	return Inventory{
		equipInventory: i.equipInventory,
		useInventory:   i.useInventory,
		setupInventory: ii,
		etcInventory:   i.etcInventory,
		cashInventory:  i.cashInventory,
	}
}

func (i Inventory) SetEtcInventory(ii ItemInventory) Inventory {
	return Inventory{
		equipInventory: i.equipInventory,
		useInventory:   i.useInventory,
		setupInventory: i.setupInventory,
		etcInventory:   ii,
		cashInventory:  i.cashInventory,
	}
}

func (i Inventory) SetCashInventory(ii ItemInventory) Inventory {
	return Inventory{
		equipInventory: i.equipInventory,
		useInventory:   i.useInventory,
		setupInventory: i.setupInventory,
		etcInventory:   i.etcInventory,
		cashInventory:  ii,
	}
}

func (i Inventory) EquipInventory() EquipInventory {
	return i.equipInventory
}

func (i Inventory) UseInventory() ItemInventory {
	return i.useInventory
}

func (i Inventory) SetupInventory() ItemInventory {
	return i.setupInventory
}

func (i Inventory) EtcInventory() ItemInventory {
	return i.etcInventory
}

func (i Inventory) CashInventory() ItemInventory {
	return i.cashInventory
}

type ItemInventory struct {
	capacity byte
	items    []Item
}

func NewItemInventory(capacity byte, items []Item) ItemInventory {
	return ItemInventory{
		capacity: capacity,
		items:    items,
	}
}

func EmptyItemInventory() ItemInventory {
	return ItemInventory{
		capacity: 4,
		items:    make([]Item, 0),
	}
}

func (i ItemInventory) Capacity() byte {
	return i.capacity
}

func (i ItemInventory) Items() []Item {
	return i.items
}

type EquipInventory struct {
	capacity byte
	items    []EquippedItem
}

func NewEquipInventory(capacity byte, items []EquippedItem) EquipInventory {
	return EquipInventory{capacity: capacity, items: items}
}

func EmptyEquipInventory() EquipInventory {
	return EquipInventory{
		capacity: 4,
		items:    make([]EquippedItem, 0),
	}
}

func (e EquipInventory) Capacity() byte {
	return e.capacity
}

func (e EquipInventory) Items() []EquippedItem {
	return e.items
}

type Item struct {
	itemId   uint32
	slot     int16
	quantity uint16
}

func NewItem(itemId uint32, slot int16, quantity uint16) Item {
	return Item{
		itemId:   itemId,
		slot:     slot,
		quantity: quantity,
	}
}

func (i Item) Slot() int16 {
	return i.slot
}

func (i Item) ItemId() uint32 {
	return i.itemId
}

func (i Item) Expiration() int64 {
	return -1
}

func (i Item) Quantity() uint16 {
	return i.quantity
}

func (i Item) Owner() string {
	return ""
}

func (i Item) Flag() uint16 {
	return 0
}

type EquippedItem struct {
	itemId        uint32
	slot          int16
	strength      uint16
	dexterity     uint16
	intelligence  uint16
	luck          uint16
	hp            uint16
	mp            uint16
	weaponAttack  uint16
	magicAttack   uint16
	weaponDefense uint16
	magicDefense  uint16
	accuracy      uint16
	avoidability  uint16
	hands         uint16
	speed         uint16
	jump          uint16
	slots         byte
}

func (i EquippedItem) NotInWeaponSlot() bool {
	if i.slot != -111 {
		return true
	}
	return false
}

func (i EquippedItem) InvertSlot() EquippedItem {
	return Clone(i).SetSlot(i.Slot() * -1).Build()
}

func (i EquippedItem) Slot() int16 {
	return i.slot
}

func (i EquippedItem) ItemId() uint32 {
	return i.itemId
}

func (i EquippedItem) InWeaponSlot() bool {
	if i.slot == -111 {
		return true
	}
	return false
}

func (i EquippedItem) IsRegularEquipment() bool {
	return i.slot > -100
}

func (i EquippedItem) Expiration() int64 {
	return -1
}

func (i EquippedItem) Slots() byte {
	return i.slots
}

func (i EquippedItem) Level() byte {
	return 0
}

func (i EquippedItem) Strength() uint16 {
	return i.strength
}

func (i EquippedItem) Dexterity() uint16 {
	return i.dexterity
}

func (i EquippedItem) Intelligence() uint16 {
	return i.intelligence
}

func (i EquippedItem) Luck() uint16 {
	return i.luck
}

func (i EquippedItem) Hp() uint16 {
	return i.hp
}

func (i EquippedItem) Mp() uint16 {
	return i.mp
}

func (i EquippedItem) WeaponAttack() uint16 {
	return i.weaponAttack
}

func (i EquippedItem) MagicAttack() uint16 {
	return i.magicAttack
}

func (i EquippedItem) WeaponDefense() uint16 {
	return i.weaponDefense
}

func (i EquippedItem) MagicDefense() uint16 {
	return i.magicDefense
}

func (i EquippedItem) Accuracy() uint16 {
	return i.accuracy
}

func (i EquippedItem) Avoidability() uint16 {
	return i.avoidability
}

func (i EquippedItem) Hands() uint16 {
	return i.hands
}

func (i EquippedItem) Speed() uint16 {
	return i.speed
}

func (i EquippedItem) Jump() uint16 {
	return i.jump
}

func (i EquippedItem) OwnerName() string {
	return ""
}

func (i EquippedItem) Flags() uint16 {
	return 0
}

func (i EquippedItem) IsEquippedCashItem() bool {
	return i.slot <= -100
}

func (i EquippedItem) Quantity() uint16 {
	return 1
}

type equippedItemBuilder struct {
	itemId        uint32
	slot          int16
	strength      uint16
	dexterity     uint16
	intelligence  uint16
	luck          uint16
	hp            uint16
	mp            uint16
	weaponAttack  uint16
	magicAttack   uint16
	weaponDefense uint16
	magicDefense  uint16
	accuracy      uint16
	avoidability  uint16
	hands         uint16
	speed         uint16
	jump          uint16
	slots         byte
}

func NewEquippedItemBuilder() *equippedItemBuilder {
	return &equippedItemBuilder{}
}

func Clone(o EquippedItem) *equippedItemBuilder {
	return &equippedItemBuilder{
		itemId:        o.itemId,
		slot:          o.slot,
		strength:      o.strength,
		dexterity:     o.dexterity,
		intelligence:  o.intelligence,
		luck:          o.luck,
		hp:            o.hp,
		mp:            o.mp,
		weaponAttack:  o.weaponAttack,
		magicAttack:   o.magicAttack,
		weaponDefense: o.weaponDefense,
		magicDefense:  o.magicDefense,
		accuracy:      o.accuracy,
		avoidability:  o.avoidability,
		hands:         o.hands,
		speed:         o.speed,
		jump:          o.jump,
		slots:         o.slots,
	}
}

func (e *equippedItemBuilder) SetItemId(itemId uint32) *equippedItemBuilder {
	e.itemId = itemId
	return e
}

func (e *equippedItemBuilder) SetSlot(slot int16) *equippedItemBuilder {
	e.slot = slot
	return e
}

func (e *equippedItemBuilder) SetStrength(strength uint16) *equippedItemBuilder {
	e.strength = strength
	return e
}

func (e *equippedItemBuilder) SetDexterity(dexterity uint16) *equippedItemBuilder {
	e.dexterity = dexterity
	return e
}

func (e *equippedItemBuilder) SetIntelligence(intelligence uint16) *equippedItemBuilder {
	e.intelligence = intelligence
	return e
}

func (e *equippedItemBuilder) SetLuck(luck uint16) *equippedItemBuilder {
	e.luck = luck
	return e
}

func (e *equippedItemBuilder) SetHp(hp uint16) *equippedItemBuilder {
	e.hp = hp
	return e
}

func (e *equippedItemBuilder) SetMp(mp uint16) *equippedItemBuilder {
	e.mp = mp
	return e
}

func (e *equippedItemBuilder) SetWeaponAttack(weaponAttack uint16) *equippedItemBuilder {
	e.weaponAttack = weaponAttack
	return e
}

func (e *equippedItemBuilder) SetMagicAttack(magicAttack uint16) *equippedItemBuilder {
	e.magicAttack = magicAttack
	return e
}

func (e *equippedItemBuilder) SetWeaponDefense(weaponDefense uint16) *equippedItemBuilder {
	e.weaponDefense = weaponDefense
	return e
}

func (e *equippedItemBuilder) SetMagicDefense(magicDefense uint16) *equippedItemBuilder {
	e.magicDefense = magicDefense
	return e
}

func (e *equippedItemBuilder) SetAccuracy(accuracy uint16) *equippedItemBuilder {
	e.accuracy = accuracy
	return e
}

func (e *equippedItemBuilder) SetAvoidability(avoidability uint16) *equippedItemBuilder {
	e.avoidability = avoidability
	return e
}

func (e *equippedItemBuilder) SetHands(hands uint16) *equippedItemBuilder {
	e.hands = hands
	return e
}

func (e *equippedItemBuilder) SetSpeed(speed uint16) *equippedItemBuilder {
	e.speed = speed
	return e
}

func (e *equippedItemBuilder) SetJump(jump uint16) *equippedItemBuilder {
	e.jump = jump
	return e
}

func (e *equippedItemBuilder) SetSlots(slots byte) *equippedItemBuilder {
	e.slots = slots
	return e
}

func (e *equippedItemBuilder) Build() EquippedItem {
	return EquippedItem{
		itemId:        e.itemId,
		slot:          e.slot,
		strength:      e.strength,
		dexterity:     e.dexterity,
		intelligence:  e.intelligence,
		luck:          e.luck,
		hp:            e.hp,
		mp:            e.mp,
		weaponAttack:  e.weaponAttack,
		magicAttack:   e.magicAttack,
		weaponDefense: e.weaponDefense,
		magicDefense:  e.magicDefense,
		accuracy:      e.accuracy,
		avoidability:  e.avoidability,
		hands:         e.hands,
		speed:         e.speed,
		jump:          e.jump,
		slots:         e.slots,
	}
}
