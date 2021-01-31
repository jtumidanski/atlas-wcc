package producers

import (
   "atlas-wcc/rest/requests"
   "context"
   "encoding/json"
   "github.com/segmentio/kafka-go"
   "log"
   "os"
   "time"
)

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
   t := requests.NewTopic(m.l)
   td, err := t.GetTopic("TOPIC_CHANNEL_SERVICE")
   if err != nil {
      m.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
   }

   w := &kafka.Writer{
      Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
      Topic:        td.Attributes.Name,
      Balancer:     &kafka.LeastBytes{},
      BatchTimeout: 50 * time.Millisecond,
   }

   e := &ChannelServerEvent{
      WorldId:   worldId,
      ChannelId: channelId,
      IpAddress: ipAddress,
      Port:      port,
      Status:    status,
   }
   r, err := json.Marshal(e)
   if err != nil {
      m.l.Fatal("[ERROR] Unable to marshall event.")
   }

   err = w.WriteMessages(context.Background(), kafka.Message{
      Key:   createKey(int(worldId)*1000 + int(channelId)),
      Value: r,
   })
   if err != nil {
      m.l.Fatal("[ERROR] Unable to produce event.")
   }
}

type ChannelServerEvent struct {
   WorldId   byte   `json:"worldId"`
   ChannelId byte   `json:"channelId"`
   IpAddress string `json:"ipAddress"`
   Port      uint32 `json:"port"`
   Status    string `json:"status"`
}
