package handler

import (
	"github.com/jtumidanski/atlas-socket/request"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

const OpInnerPortal uint16 = 0x65
const InnerPortal = "inner_portal"

func InnerPortalHandlerProducer(l logrus.FieldLogger) Producer {
	return func() (uint16, request.Handler) {
		return OpInnerPortal, SpanHandlerDecorator(l, InnerPortal, func(l logrus.FieldLogger, span opentracing.Span) request.Handler {
			return ValidatorHandler(LoggedInValidator(l, span), NoOpHandler)
		})
	}
}
