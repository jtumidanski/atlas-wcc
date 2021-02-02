package consumers

import (
   "atlas-wcc/socket/response/writer"
   "log"
)

type itemPickedUpEvent struct {
   CharacterId uint32 `json:"characterId"`
   ItemId uint32 `json:"itemId"`
   Quantity uint32 `json:"quantity"`
}

func ItemPickedUpEventCreator() EmptyEventCreator {
   return func() interface{} {
      return &itemPickedUpEvent{}
   }
}

func HandleItemPickedUpEvent() ChannelEventProcessor {
   return func(l *log.Logger, wid byte, cid byte, e interface{}) {
      if event, ok := e.(*itemPickedUpEvent); ok {
         as := getSessionForCharacterId(event.CharacterId)
         if as == nil {
            l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
            return
         }
         (*as).Announce(writer.WriteShowItemGain(event.ItemId, event.Quantity))
         (*as).Announce(writer.WriteEnableActions())
      } else {
         l.Printf("[ERROR] unable to cast event provided to handler [HandleItemPickedUpEvent]")
      }
   }
}
