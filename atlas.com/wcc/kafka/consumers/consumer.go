package consumers

import (
   "atlas-wcc/rest/requests"
   "context"
   "encoding/json"
   "fmt"
   "github.com/segmentio/kafka-go"
   "log"
   "os"
   "time"
)

type Consumer struct {
   l         *log.Logger
   ctx       context.Context
   worldId   byte
   channelId byte
   h         EventHandler
}

func NewConsumer(l *log.Logger, ctx context.Context, worldId byte, channelId byte, h EventHandler) Consumer {
   return Consumer{l, ctx, worldId, channelId, h}
}

type EventHandler interface {
   topicToken() string

   emptyEventCreator() interface{}

   eventProcessor(*log.Logger, interface{})
}

func (c Consumer) Init() {
   t := requests.NewTopic(c.l)
   td, err := t.GetTopic(c.h.topicToken())
   if err != nil {
      c.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
   }

   c.l.Printf("[INFO] creating topic consumer for %s", td.Attributes.Name)
   r := kafka.NewReader(kafka.ReaderConfig{
      Brokers: []string{os.Getenv("BOOTSTRAP_SERVERS")},
      Topic:   td.Attributes.Name,
      GroupID: fmt.Sprintf("World Channel Coordinator %d %d", c.worldId, c.channelId),
      MaxWait: 50 * time.Millisecond,
   })
   for {
      msg, err := r.ReadMessage(c.ctx)
      if err != nil {
         panic("Could not successfully read message " + err.Error())
      }

      event := c.h.emptyEventCreator()
      err = json.Unmarshal(msg.Value, &event)
      if err != nil {
         c.l.Println("Could not unmarshal event into event class ", msg.Value)
      } else {
         c.h.eventProcessor(c.l, event)
      }
   }
}
