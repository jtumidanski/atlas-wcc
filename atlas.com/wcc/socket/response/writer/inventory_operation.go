package writer

import (
	"atlas-wcc/domain"
	"atlas-wcc/socket/response"
	"fmt"
)

const OpCodeInventoryOperation uint16 = 0x1D

type Modification struct {
	Mode          byte
	ItemId        uint32
	InventoryType byte
	Quantity      uint16
	Position      int16
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
			w.WriteInt16(mod.Position + 1)
		}
		switch mod.Mode {
		case 0:
			addItem(w, mod)
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
	if (mod.Position + 1) < 0 {
		return 2
	}
	return movement
}

func moveItem(w *response.Writer, movement int8, mod Modification) int8 {
	w.WriteInt16(mod.Position + 1)
	if (mod.Position+1) < 0 || (mod.OldPosition+1) < 0 {
		if (mod.OldPosition + 1) < 0 {
			return 1
		}
		return 2
	}
	return movement
}

func updateQuantity(w *response.Writer, mod Modification) {
	w.WriteShort(mod.Quantity)
}

func addItem(w *response.Writer, mod Modification) {
	addItemInfoZero(w, domain.NewItem(mod.ItemId, int8(mod.Position+1), mod.Quantity), true)
}
