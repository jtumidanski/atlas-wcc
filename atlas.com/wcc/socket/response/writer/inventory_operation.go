package writer

import (
	"atlas-wcc/domain"
	"atlas-wcc/socket/response"
	"fmt"
)

const OpCodeInventoryOperation uint16 = 0x1D

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
	InventoryType byte
	Item          InventoryItem
	OldPosition   int16
}

type ModifyInventory struct {
	UpdateTick    bool
	Modifications []Modification
}

func WriteCharacterInventoryModification(input ModifyInventory) []byte {
	w := response.NewWriter()
	w.WriteShort(OpCodeInventoryOperation)
	w.WriteBool(input.UpdateTick)
	w.WriteByte(byte(len(input.Modifications)))
	addMovement := int8(-1)
	for _, mod := range input.Modifications {
		w.WriteByte(mod.Mode)
		w.WriteByte(mod.InventoryType)
		if mod.Mode == 2 {
			w.WriteInt16(mod.OldPosition + 1)
		} else {
			w.WriteInt16(mod.Item.Slot() + 1)
		}
		switch mod.Mode {
		case 0:
			if mod.InventoryType == 1 {
				if val, ok := mod.Item.(*domain.EquippedItem); ok {
					addEquipmentInfoZero(w, *val, true)
				}
			} else {
				if val, ok := mod.Item.(*domain.Item); ok {
					addItemInfoZero(w, *val, true)
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

func removeItem(movement int8, mod Modification) int8 {
	if (mod.Item.Slot() + 1) < 0 {
		return 2
	}
	return movement
}

func moveItem(w *response.Writer, movement int8, mod Modification) int8 {
	w.WriteInt16(mod.Item.Slot())
	if (mod.Item.Slot() + 1) < 0 || (mod.OldPosition + 1) < 0 {
		if (mod.OldPosition + 1) < 0 {
			return 1
		}
		return 2
	}
	return movement
}

func updateQuantity(w *response.Writer, mod Modification) {
	w.WriteShort(mod.Item.Quantity())
}
