package handler

import (
	"atlas-wcc/cashshop"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const (
	OpEnterCashShop     = 0x28
	OpTouchingCashShop  = 0xE4
	OpCashShopOperation = 0xE5
)

func EnterCashShopHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, _ *request.RequestReader) {
		cashshop.RequestCashShopEntry(l, span)(worldId, channelId, s.CharacterId())
	}
}

func TouchingCashShopHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		cashshop.UpdateCashAmounts(l, span)(worldId, channelId, s.CharacterId())
	}
}

type itemPurchaseRequest struct {
	cashIndex    uint32
	serialNumber uint32
}

func (r itemPurchaseRequest) SerialNumber() uint32 {
	return r.serialNumber
}

func (r itemPurchaseRequest) CashIndex() uint32 {
	return r.cashIndex
}

type packagePurchaseRequest struct {
	cashIndex    uint32
	serialNumber uint32
}

type giftPurchaseRequest struct {
	birthday     uint32
	serialNumber uint32
	recipient    string
	message      string
}

type modifyWishlistRequest struct {
	serialNumbers []uint32
}

func (r modifyWishlistRequest) SerialNumbers() []uint32 {
	return r.serialNumbers
}

type adjustSlotsViaButtonRequest struct {
	cashIndex uint32
	slotType  byte
}

type adjustSlotsViaItemRequest struct {
	cashIndex    uint32
	serialNumber uint32
}

type adjustStorageSlotsViaButtonRequest struct {
	cashIndex uint32
}

type adjustStorageSlotsViaItemRequest struct {
	cashIndex    uint32
	serialNumber uint32
}

type increaseCharacterSlotsRequest struct {
	cashIndex    uint32
	serialNumber uint32
}

type takeFromCashInventoryRequest struct {
	cashId uint32
}

type putIntoCashInventoryRequest struct {
	cashId        uint32
	inventoryType byte
}

type crushRingRequest struct {
	birthday     uint32
	cashIndex    uint32
	serialNumber uint32
	recipient    string
	message      string
}

type mesoItemRequest struct {
	serialNumber uint32
}

type friendshipRingRequest struct {
	birthday     uint32
	cashIndex    uint32
	serialNumber uint32
	recipient    string
	message      string
}

type nameChangeRequest struct {
	serialNumber uint32
	oldName      string
	newName      string
}

type worldTransferRequest struct {
	serialNumber uint32
	newWorldId   uint32
}

func readCashShopOperation(r *request.RequestReader) interface{} {
	action := r.ReadByte()
	if action == 0x03 {
		_ = r.ReadByte()
		cashIndex := r.ReadUint32()
		serialNumber := r.ReadUint32()
		return &itemPurchaseRequest{cashIndex: cashIndex, serialNumber: serialNumber}
	}
	if action == 0x1E {
		_ = r.ReadByte()
		cashIndex := r.ReadUint32()
		serialNumber := r.ReadUint32()
		return &packagePurchaseRequest{cashIndex: cashIndex, serialNumber: serialNumber}
	}
	if action == 0x04 {
		birthday := r.ReadUint32()
		serialNumber := r.ReadUint32()
		recipient := r.ReadAsciiString()
		message := r.ReadAsciiString()
		return &giftPurchaseRequest{
			birthday:     birthday,
			serialNumber: serialNumber,
			recipient:    recipient,
			message:      message,
		}
	}
	if action == 0x05 {
		serialNumbers := make([]uint32, 0)
		for i := 0; i < 10; i++ {
			serialNumbers = append(serialNumbers, r.ReadUint32())
		}
		return &modifyWishlistRequest{serialNumbers: serialNumbers}
	}
	if action == 0x06 {
		r.Skip(1)
		cashIndex := r.ReadUint32()
		mode := r.ReadByte()
		if mode == 0 {
			slotType := r.ReadByte()
			return &adjustSlotsViaButtonRequest{cashIndex: cashIndex, slotType: slotType}
		}
		serialNumber := r.ReadUint32()
		return &adjustSlotsViaItemRequest{cashIndex: cashIndex, serialNumber: serialNumber}
	}
	if action == 0x07 {
		r.Skip(1)
		cashIndex := r.ReadUint32()
		mode := r.ReadByte()
		if mode == 0 {
			return adjustStorageSlotsViaButtonRequest{cashIndex: cashIndex}
		}
		serialNumber := r.ReadUint32()
		return adjustStorageSlotsViaItemRequest{cashIndex: cashIndex, serialNumber: serialNumber}
	}
	if action == 0x08 {
		r.Skip(1)
		cashIndex := r.ReadUint32()
		serialNumber := r.ReadUint32()
		return &increaseCharacterSlotsRequest{cashIndex: cashIndex, serialNumber: serialNumber}
	}
	if action == 0x0D {
		cashId := r.ReadUint32()
		return &takeFromCashInventoryRequest{cashId: cashId}
	}
	if action == 0x0E {
		cashId := r.ReadUint32()
		r.Skip(4)
		inventoryType := r.ReadByte()
		return &putIntoCashInventoryRequest{cashId: cashId, inventoryType: inventoryType}
	}
	if action == 0x1D {
		birthday := r.ReadUint32()
		cashIndex := r.ReadUint32()
		serialNumber := r.ReadUint32()
		recipient := r.ReadAsciiString()
		message := r.ReadAsciiString()
		return &crushRingRequest{
			birthday:     birthday,
			cashIndex:    cashIndex,
			serialNumber: serialNumber,
			recipient:    recipient,
			message:      message,
		}
	}
	if action == 0x20 {
		serialNumber := r.ReadUint32()
		return &mesoItemRequest{serialNumber: serialNumber}
	}
	if action == 0x23 {
		birthday := r.ReadUint32()
		cashIndex := r.ReadUint32()
		serialNumber := r.ReadUint32()
		recipient := r.ReadAsciiString()
		message := r.ReadAsciiString()
		return &friendshipRingRequest{
			birthday:     birthday,
			cashIndex:    cashIndex,
			serialNumber: serialNumber,
			recipient:    recipient,
			message:      message,
		}
	}
	if action == 0x2E {
		serialNumber := r.ReadUint32()
		oldName := r.ReadAsciiString()
		newName := r.ReadAsciiString()
		return &nameChangeRequest{
			serialNumber: serialNumber,
			oldName:      oldName,
			newName:      newName,
		}
	}
	if action == 0x31 {
		serialNumber := r.ReadUint32()
		newWorldId := r.ReadUint32()
		return &worldTransferRequest{serialNumber: serialNumber, newWorldId: newWorldId}
	}
	return nil
}

func CashShopOperationHandler(l logrus.FieldLogger, span opentracing.Span, _ byte, _ byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readCashShopOperation(r)
		if val, ok := p.(*itemPurchaseRequest); ok {
			cashshop.RequestItemPurchase(l, span)(s.CharacterId(), val.CashIndex(), val.SerialNumber())
		}
		if val, ok := p.(*modifyWishlistRequest); ok {
			cashshop.ModifyWishlist(l, span)(s.CharacterId(), val.SerialNumbers())
			return
		}
	}
}
