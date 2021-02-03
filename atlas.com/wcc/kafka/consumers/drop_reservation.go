package consumers

import (
   "atlas-wcc/socket/response/writer"
   "log"
)

type dropReservationEvent struct {
   CharacterId uint32 `json:"characterId"`
   DropId uint32 `json:"dropId"`
   Type string `json:"type"`
}

func DropReservationEventCreator() EmptyEventCreator {
   return func() interface{} {
      return &dropReservationEvent{}
   }
}

func HandleDropReservationEvent() ChannelEventProcessor {
   return func(l *log.Logger, wid byte, cid byte, e interface{}) {
      if event, ok := e.(*dropReservationEvent); ok {
         if event.Type == "SUCCESS" {
            return
         }

         as := getSessionForCharacterId(event.CharacterId)
         if as == nil {
            l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
            return
         }
         (*as).Announce(writer.WriteEnableActions())
      } else {
         l.Printf("[ERROR] unable to cast event provided to handler [HandleDropReservationEvent]")
      }
   }
}
