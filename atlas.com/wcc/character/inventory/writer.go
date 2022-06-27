package inventory

import (
	"atlas-wcc/socket/response"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

const OpCodeShowStatusInfo uint16 = 0x27
const OpCodeInventoryOperation uint16 = 0x1D

func WriteShowInventoryFull(l logrus.FieldLogger) []byte {
	w := response.NewWriter(l)
	w.WriteShort(OpCodeShowStatusInfo)
	w.WriteByte(0)
	w.WriteByte(0xFF)
	w.WriteInt(0)
	w.WriteInt(0)
	return w.Bytes()
}

type InventoryItem interface {
	Slotter
	Quantity
}

type Slotter interface {
	Slot() int16
}

type Quantity interface {
	Quantity() uint16
}

type Modification struct {
	Mode          byte
	InventoryType int8
	Item          InventoryItem
	OldPosition   int16
}

type ModifyInventory struct {
	UpdateTick    bool
	Modifications []Modification
}

func WriteCharacterInventoryModification(l logrus.FieldLogger) func(input ModifyInventory) []byte {
	return func(input ModifyInventory) []byte {
		w := response.NewWriter(l)
		w.WriteShort(OpCodeInventoryOperation)
		w.WriteBool(input.UpdateTick)
		w.WriteByte(byte(len(input.Modifications)))
		addMovement := int8(-1)
		for _, mod := range input.Modifications {
			w.WriteByte(mod.Mode)
			w.WriteInt8(mod.InventoryType)
			if mod.Mode == 2 {
				w.WriteInt16(mod.OldPosition)
			} else {
				w.WriteInt16(mod.Item.Slot())
			}
			switch mod.Mode {
			case 0:
				if mod.InventoryType == 1 {
					if val, ok := mod.Item.(*EquippedItem); ok {
						addEquipmentInfoZero(w, *val, true)
					}
				} else {
					if val, ok := mod.Item.(Item); ok {
						addItemInfoZero(w, val, true)
					}
				}
				break
			case 1:
				updateQuantity(w, mod)
				break
			case 2:
				addMovement = moveItem(w, addMovement, mod)
				break
			case 3:
				addMovement = removeItem(addMovement, mod)
				break
			default:
				panic(fmt.Sprintf("unsupported inventory mode %d", mod.Mode))
			}
		}
		if addMovement > -1 {
			w.WriteInt8(addMovement)
		}
		return w.Bytes()
	}
}

func removeItem(movement int8, mod Modification) int8 {
	slot := mod.Item.Slot()
	if (slot) < 0 {
		return 2
	}
	return movement
}

func moveItem(w *response.Writer, movement int8, mod Modification) int8 {
	slot := mod.Item.Slot()
	oldSlot := mod.OldPosition

	w.WriteInt16(slot)
	if slot < 0 || oldSlot < 0 {
		if oldSlot < 0 {
			return 1
		}
		return 2
	}
	return movement
}

func updateQuantity(w *response.Writer, mod Modification) {
	w.WriteShort(mod.Item.Quantity())
}

func AddInventoryInfo(w *response.Writer, equipment []EquippedItem, inventory Inventory) {
	equip := inventory.EquipInventory()
	use := inventory.UseInventory()
	setup := inventory.SetupInventory()
	etc := inventory.EtcInventory()
	cash := inventory.CashInventory()
	w.WriteByte(equip.Capacity())
	w.WriteByte(use.Capacity())
	w.WriteByte(setup.Capacity())
	w.WriteByte(etc.Capacity())
	w.WriteByte(cash.Capacity())

	w.WriteLong(uint64(getTime(-2)))

	for _, e := range equipment {
		if e.IsRegularEquipment() {
			addEquipmentInfo(w, e)
		}
	}

	w.WriteShort(0)
	for _, e := range equipment {
		if e.IsEquippedCashItem() {
			addEquipmentInfo(w, e)
		}
	}

	w.WriteShort(0)
	for _, e := range equip.Items() {
		addEquipmentInfo(w, e)
	}

	w.WriteInt(0)
	for _, i := range use.Items() {
		addItemInfo(w, i)
	}

	w.WriteByte(0)
	for _, i := range setup.Items() {
		addItemInfo(w, i)
	}

	w.WriteByte(0)
	for _, i := range etc.Items() {
		addItemInfo(w, i)
	}

	w.WriteByte(0)
	for _, i := range cash.Items() {
		addItemInfo(w, i)
	}
}

func addItemInfo(w *response.Writer, i Item) {
	addItemInfoZero(w, i, false)
}

func addItemInfoZero(w *response.Writer, i Item, zeroPosition bool) {
	if !zeroPosition {
		w.WriteInt8(int8(i.Slot()))
	}
	w.WriteByte(2)
	w.WriteInt(i.ItemId())
	w.WriteBool(false)
	w.WriteLong(uint64(getTime(i.Expiration())))
	w.WriteShort(i.Quantity())
	w.WriteAsciiString(i.Owner())
	w.WriteShort(i.Flag())
}

func addEquipmentInfo(w *response.Writer, e EquippedItem) {
	addEquipmentInfoZero(w, e, false)
}

func addEquipmentInfoZero(w *response.Writer, e EquippedItem, zeroPosition bool) {
	slot := e.Slot()
	if !zeroPosition {
		if slot < 0 {
			slot *= -1
		}
		if slot > 100 {
			w.WriteShort(uint16(slot - 100))
		} else {
			w.WriteShort(uint16(slot))
		}
	}

	w.WriteByte(1)
	w.WriteInt(e.ItemId())
	w.WriteBool(false)
	w.WriteLong(uint64(getTime(e.Expiration())))
	w.WriteByte(e.Slots())
	w.WriteByte(e.Level())
	w.WriteShort(e.Strength())
	w.WriteShort(e.Dexterity())
	w.WriteShort(e.Intelligence())
	w.WriteShort(e.Luck())
	w.WriteShort(e.Hp())
	w.WriteShort(e.Mp())
	w.WriteShort(e.WeaponAttack())
	w.WriteShort(e.MagicAttack())
	w.WriteShort(e.WeaponDefense())
	w.WriteShort(e.MagicDefense())
	w.WriteShort(e.Accuracy())
	w.WriteShort(e.Avoidability())
	w.WriteShort(e.Hands())
	w.WriteShort(e.Speed())
	w.WriteShort(e.Jump())
	w.WriteAsciiString(e.OwnerName())
	w.WriteShort(e.Flags())

	w.WriteByte(0)
	w.WriteByte(0)
	w.WriteInt(0)
	w.WriteInt(0)
	w.WriteLong(0)
	w.WriteLong(uint64(getTime(-2)))
	w.WriteInt32(-1)
}

const (
	DefaultTime int64 = 150842304000000000
	ZeroTime    int64 = 94354848000000000
	Permanent   int64 = 150841440000000000
)

func getTime(utcTimestamp int64) int64 {
	if utcTimestamp < 0 && utcTimestamp >= -3 {
		if utcTimestamp == -1 {
			return DefaultTime //high number ll
		} else if utcTimestamp == -2 {
			return ZeroTime
		} else {
			return Permanent
		}
	}

	ftUtOffset := 116444736010800000 + (10000 * timeNow())
	return utcTimestamp*10000 + ftUtOffset
}

func timeNow() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func AddCashItemInformation(w *response.Writer, index int, item Item, accountId uint32) {
	// TODO determine if it's a gift, a ring, or a pet. All assumptions = no
	w.WriteLong(uint64(index))
	w.WriteInt(accountId)
	w.WriteInt(0)
	w.WriteInt(item.ItemId())
	w.WriteInt(item.ItemId()) // needs to be SN
	w.WriteShort(item.Quantity())
	addPaddedGiftOriginatorName(w, "") //gift from
	w.WriteLong(uint64(getTime(item.Expiration())))
	w.WriteLong(0)
}

func addPaddedGiftOriginatorName(w *response.Writer, name string) {
	if len(name) > 13 {
		name = name[:13]
	}
	padSize := 13 - len(name)
	w.WriteByteArray([]byte(name))
	for i := 0; i < padSize; i++ {
		w.WriteByte(0x0)
	}
}
