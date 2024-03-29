package handler

import (
	"atlas-wcc/character"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpMoveItem uint16 = 0x47
const MoveItem = "move_item"

func MoveItemHandlerProducer(l logrus.FieldLogger, worldId byte, channelId byte) Producer {
	return func() (uint16, request.Handler) {
		return OpMoveItem, SpanHandlerDecorator(l, MoveItem, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(LoggedInValidator(l, span), MoveItemHandler(l, span, worldId, channelId))
		})
	}
}

type moveItemRequest struct {
	inventoryType int8
	source        int16
	action        int16
	quantity      int16
}

func (r moveItemRequest) Source() int16 {
	return r.source
}

func (r moveItemRequest) Action() int16 {
	return r.action
}

func (r moveItemRequest) InventoryType() int8 {
	return r.inventoryType
}

func (r moveItemRequest) Quantity() int16 {
	return r.quantity
}

func readMoveItemRequest(reader *request.RequestReader) moveItemRequest {
	reader.Skip(4)
	inventoryType := reader.ReadInt8()
	source := reader.ReadInt16()
	action := reader.ReadInt16()
	quantity := reader.ReadInt16()
	return moveItemRequest{inventoryType: inventoryType, source: source, action: action, quantity: quantity}
}

func MoveItemHandler(l logrus.FieldLogger, span opentracing.Span, worldId byte, channelId byte) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readMoveItemRequest(r)
		// adjust for client indexing positive from 1 not 0
		source := p.Source()
		action := p.Action()

		if p.Source() < 0 && p.Action() > 0 {
			character.UnequipItem(l, span)(s.CharacterId(), source, action)
		} else if p.Action() < 0 {
			character.EquipItem(l, span)(s.CharacterId(), source, action)
		} else if p.Action() == 0 {
			character.DropItem(l, span)(worldId, channelId, s.CharacterId(), p.InventoryType(), p.Source(), p.Quantity())
		} else {
			character.MoveItem(l, span)(s.CharacterId(), p.InventoryType(), p.Source(), p.Action())
		}
	}
}
