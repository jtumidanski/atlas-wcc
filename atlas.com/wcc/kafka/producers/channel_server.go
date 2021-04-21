package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type channelServerEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	IpAddress string `json:"ipAddress"`
	Port      uint32 `json:"port"`
	Status    string `json:"status"`
}

var ChannelServer = func(l logrus.FieldLogger, ctx context.Context) *channelServer {
	return &channelServer{
		l:   l,
		ctx: ctx,
	}
}

type channelServer struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func (m *channelServer) Start(worldId byte, channelId byte, ipAddress string, port uint32) {
	m.emit(worldId, channelId, ipAddress, port, "STARTED")
}

func (m *channelServer) Shutdown(worldId byte, channelId byte, ipAddress string, port uint32) {
	m.emit(worldId, channelId, ipAddress, port, "SHUTDOWN")
}

func (m *channelServer) emit(worldId byte, channelId byte, ipAddress string, port uint32, status string) {
	e := &channelServerEvent{
		WorldId:   worldId,
		ChannelId: channelId,
		IpAddress: ipAddress,
		Port:      port,
		Status:    status,
	}
	produceEvent(m.l, "TOPIC_CHANNEL_SERVICE", createKey(int(worldId)*1000+int(channelId)), e)
}
