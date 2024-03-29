package handler

import (
	"atlas-wcc/character/keymap"
	"atlas-wcc/session"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpChangeKeyMap uint16 = 0x87
const ChangeKeyMap = "change_key_map"

func ChangeKeyMapHandlerProducer(l logrus.FieldLogger) Producer {
	return func() (uint16, request.Handler) {
		return OpChangeKeyMap, SpanHandlerDecorator(l, ChangeKeyMap, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(LoggedInValidator(l, span), ChangeKeyMapHandler(l, span))
		})
	}
}

type changeKeyMapRequest struct {
	available bool
	changes   []change
}

type change struct {
	key        int32
	changeType int8
	action     int32
}

func (c change) Key() int32 {
	return c.key
}

func (c change) Type() int8 {
	return c.changeType
}

func (c change) Action() int32 {
	return c.action
}

func readChangeKeyMapRequest(reader *request.RequestReader) interface{} {
	available := len(reader.GetBuffer()) >= 8
	mode := int32(-1)
	if available {
		mode = reader.ReadInt32()
		if mode == 0 {
			changeCount := reader.ReadInt32()
			changes := make([]change, 0)
			for i := int32(0); i < changeCount; i++ {
				key := reader.ReadInt32()
				changeType := reader.ReadInt8()
				action := reader.ReadInt32()
				changes = append(changes, change{key: key, changeType: changeType, action: action})
			}
			return changeKeyMapRequest{available: available, changes: changes}
		}
	}
	return nil
}

func ChangeKeyMapHandler(l logrus.FieldLogger, span opentracing.Span) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readChangeKeyMapRequest(r)
		if packet, ok := p.(changeKeyMapRequest); ok {
			if packet.available {
				changes := make([]keymap.Change, 0)
				for _, c := range packet.changes {
					changes = append(changes, keymap.Change{Key: c.Key(), ChangeType: c.Type(), Action: c.Action()})
				}
				keymap.ChangeKeyMap(l, span)(s.CharacterId(), changes)
			}
		}
	}
}
