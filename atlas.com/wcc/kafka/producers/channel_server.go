package producers

import (
	"context"
	"log"
)

type ChannelServerEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	IpAddress string `json:"ipAddress"`
	Port      uint32 `json:"port"`
	Status    string `json:"status"`
}

type ChannelServer struct {
	l   *log.Logger
	ctx context.Context
}

func NewChannelServer(l *log.Logger, ctx context.Context) *ChannelServer {
	return &ChannelServer{l, ctx}
}

func (m *ChannelServer) EmitStart(worldId byte, channelId byte, ipAddress string, port uint32) {
	m.emit(worldId, channelId, ipAddress, port, "STARTED")
}

func (m *ChannelServer) EmitShutdown(worldId byte, channelId byte, ipAddress string, port uint32) {
	m.emit(worldId, channelId, ipAddress, port, "SHUTDOWN")
}

func (m *ChannelServer) emit(worldId byte, channelId byte, ipAddress string, port uint32, status string) {
	e := &ChannelServerEvent{
		WorldId:   worldId,
		ChannelId: channelId,
		IpAddress: ipAddress,
		Port:      port,
		Status:    status,
	}
	ProduceEvent(m.l, "TOPIC_CHANNEL_SERVICE", createKey(int(worldId)*1000+int(channelId)), e)
}
