package movement

import (
	"atlas-wcc/session"
	"atlas-wcc/socket/request/handler"
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpNpcAction uint16 = 0xC5
const NPCAction = "npc_action"

func HandleNPCActionProducer(l logrus.FieldLogger) handler.Producer {
	return func() (uint16, request.Handler) {
		return OpNpcAction, handler.SpanHandlerDecorator(l, NPCAction, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return handler.ValidatorHandler(handler.LoggedInValidator(l, span), handleNPCAction(l))
		})
	}
}

func handleNPCAction(l logrus.FieldLogger) func(s session.Model, r *request.RequestReader) {
	return func(s session.Model, r *request.RequestReader) {
		p := readNPCAction(r)
		if val, ok := p.(*npcAnimationRequest); ok {
			err := session.Announce(s, writeNPCAnimation(l)(val.ObjectId(), val.Second(), val.Third()))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return
		}
		if val, ok := p.(*npcMoveRequest); ok {
			err := session.Announce(s, writeNPCMove(l)(val.Movement()))
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
			return
		}
		l.Warnf("Received a unhandled [NPCActionRequest]")
	}
}
