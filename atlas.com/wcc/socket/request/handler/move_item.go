package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/session"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/sirupsen/logrus"
)

const OpMoveItem uint16 = 0x47

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

func MoveItemHandler() request2.MessageHandler {
	return func(l logrus.FieldLogger, s *session.Model, r *request.RequestReader) {
		p := readMoveItemRequest(r)
		// adjust for client indexing positive from 1 not 0
		source := p.Source()
		action := p.Action()

		if p.Source() < 0 && p.Action() > 0 {
			producers.UnequipItem(l)(s.CharacterId(), source, action)
		} else if p.Action() < 0 {
			producers.EquipItem(l)(s.CharacterId(), source, action)
		} else if p.Action() == 0 {
			producers.DropItem(l)(s.WorldId(), s.ChannelId(), s.CharacterId(), p.InventoryType(), p.Source(), p.Quantity())
		} else {
			producers.MoveItem(l)(s.CharacterId(), p.InventoryType(), p.Source(), p.Action())
		}
	}
}
