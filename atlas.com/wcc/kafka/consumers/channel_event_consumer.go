package consumers

import "log"

type ChannelEventProcessor func(*log.Logger, byte, byte, interface{})
