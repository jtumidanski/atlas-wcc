package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type channelServerEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	IpAddress string `json:"ipAddress"`
	Port      uint32 `json:"port"`
	Status    string `json:"status"`
}

func StartChannelServer(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, ipAddress string, port uint32) {
	producer := ProduceEvent(l, span, "TOPIC_CHANNEL_SERVICE")
	return func(worldId byte, channelId byte, ipAddress string, port uint32) {
		emitChannelServer(producer, worldId, channelId, ipAddress, port, "STARTED")
	}
}

func ShutdownChannelServer(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, ipAddress string, port uint32) {
	producer := ProduceEvent(l, span, "TOPIC_CHANNEL_SERVICE")
	return func(worldId byte, channelId byte, ipAddress string, port uint32) {
		emitChannelServer(producer, worldId, channelId, ipAddress, port, "SHUTDOWN")
	}
}

func emitChannelServer(producer func(key []byte, event interface{}), worldId byte, channelId byte, ipAddress string, port uint32, status string) {
	e := &channelServerEvent{
		WorldId:   worldId,
		ChannelId: channelId,
		IpAddress: ipAddress,
		Port:      port,
		Status:    status,
	}
	producer(CreateKey(int(worldId)*1000+int(channelId)), e)
}
