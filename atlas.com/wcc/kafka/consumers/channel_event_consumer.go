package consumers

import (
	"github.com/sirupsen/logrus"
)

type ChannelEventProcessor func(logrus.FieldLogger, byte, byte, interface{})
