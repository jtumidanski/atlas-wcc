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

type PortalEnter struct {
   l   *log.Logger
   ctx context.Context
}

func NewPortalEnter(l *log.Logger, ctx context.Context) *PortalEnter {
   return &PortalEnter{l, ctx}
}

func (m *PortalEnter) EmitEnter(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
   m.emit(worldId, channelId, mapId, portalId, characterId)
}

func (m *PortalEnter) emit(worldId byte, channelId byte, mapId uint32, portalId uint32, characterId uint32) {
   t := requests.NewTopic(m.l)
   td, err := t.GetTopic("TOPIC_ENTER_PORTAL")
   if err != nil {
      m.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
   }

   w := &kafka.Writer{
      Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
      Topic:        td.Attributes.Name,
      Balancer:     &kafka.LeastBytes{},
      BatchTimeout: 50 * time.Millisecond,
   }

   e := &PortalEnterCommand{
      WorldId:     worldId,
      ChannelId:   channelId,
      MapId:       mapId,
      PortalId:    portalId,
      CharacterId: characterId,
   }
   r, err := json.Marshal(e)
   if err != nil {
      m.l.Fatal("[ERROR] Unable to marshall event.")
   }

   err = w.WriteMessages(context.Background(), kafka.Message{
      Key:   createKey(int(characterId)),
      Value: r,
   })
   if err != nil {
      m.l.Fatal("[ERROR] Unable to produce event.")
   }
}

type PortalEnterCommand struct {
   WorldId     byte   `json:"worldId"`
   ChannelId   byte   `json:"channelId"`
   MapId       uint32 `json:"mapId"`
   PortalId    uint32 `json:"portalId"`
   CharacterId uint32 `json:"characterId"`
}
