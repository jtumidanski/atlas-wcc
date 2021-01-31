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

type CharacterMovement struct {
   l   *log.Logger
   ctx context.Context
}

func NewCharacterMovement(l *log.Logger, ctx context.Context) *CharacterMovement {
   return &CharacterMovement{l, ctx}
}

func (m *CharacterMovement) EmitMovement(worldId byte, channelId byte, characterId uint32, x int16, y int16, stance byte, rawMovement []byte) {
   t := requests.NewTopic(m.l)
   td, err := t.GetTopic("TOPIC_CHARACTER_MOVEMENT")
   if err != nil {
      m.l.Fatal("[ERROR] Unable to retrieve topic for consumer.")
   }

   w := &kafka.Writer{
      Addr:         kafka.TCP(os.Getenv("BOOTSTRAP_SERVERS")),
      Topic:        td.Attributes.Name,
      Balancer:     &kafka.LeastBytes{},
      BatchTimeout: 50 * time.Millisecond,
   }

   e := &CharacterMovementEvent{
      WorldId:     worldId,
      ChannelId:   channelId,
      CharacterId: characterId,
      X:           x,
      Y:           y,
      Stance:      stance,
      RawMovement: rawMovement,
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

type CharacterMovementEvent struct {
   WorldId     byte   `json:"worldId"`
   ChannelId   byte   `json:"channelId"`
   CharacterId uint32 `json:"characterId"`
   X           int16  `json:"x"`
   Y           int16  `json:"y"`
   Stance      byte   `json:"stance"`
   RawMovement []byte `json:"rawMovement"`
}
