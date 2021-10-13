package consumers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type ChannelEventProcessor func(logrus.FieldLogger, opentracing.Span, byte, byte, interface{})
