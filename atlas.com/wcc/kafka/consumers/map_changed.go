package consumers

import (
   "atlas-wcc/rest/requests"
   "atlas-wcc/socket/response/writer"
   "context"
   "encoding/json"
   "fmt"
   "github.com/segmentio/kafka-go"
   "log"
   "os"
   "time"
)

type MapChanged struct {
   l   *log.Logger
   ctx context.Context
}

func NewMapChanged(l *log.Logger, ctx context.Context) *MapChanged {
   return &MapChanged{l, ctx}
}

func (mc *MapChanged) Init(worldId byte, channelId byte) {
   t := requests.NewTopic(mc.l)
   td, err := t.GetTopic("TOPIC_CHANGE_MAP_EVENT")
   if err != nil {
      mc.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
   }

   r := kafka.NewReader(kafka.ReaderConfig{
      Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
      Topic:   td.Attributes.Name,
      GroupID: fmt.Sprintf("World Channel Coordinator %d %d", worldId, channelId),
      MaxWait: 50 * time.Millisecond,
   })
   for {
      msg, err := r.ReadMessage(mc.ctx)
      if err != nil {
         panic("Could not successfully read message " + err.Error())
      }

      var event MapChangedEvent
      err = json.Unmarshal(msg.Value, &event)
      if err != nil {
         mc.l.Println("Could not unmarshal event into event class ", msg.Value)
      } else {
         mc.processEvent(event)
      }
   }
}

func (mc *MapChanged) processEvent(event MapChangedEvent) {
   as := getSessionForCharacterId(event.CharacterId)
   if as == nil {
      return
   }

   catt, err := requests.GetCharacterAttributesById(event.CharacterId)
   if err != nil {
      return
   }

   (*as).Announce(writer.WriteWarpToMap(event.ChannelId, event.MapId, event.PortalId, catt.Data().Attributes.Hp))
}

type MapChangedEvent struct {
   WorldId     byte   `json:"worldId"`
   ChannelId   byte   `json:"channelId"`
   MapId       uint32 `json:"mapId"`
   PortalId    uint32 `json:"portalId"`
   CharacterId uint32 `json:"characterId"`
}
