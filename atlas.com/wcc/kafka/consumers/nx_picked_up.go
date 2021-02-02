package consumers

import (
   "atlas-wcc/socket/response/writer"
   "fmt"
   "log"
)

const nxGainFormat = "You have earned #e#b%d NX#k#n."

type nxPickedUpEvent struct {
   CharacterId uint32 `json:"characterId"`
   Gain        uint32 `json:"gain"`
}

func NXPickedUpEventCreator() EmptyEventCreator {
   return func() interface{} {
      return &nxPickedUpEvent{}
   }
}

func HandleNXPickedUpEvent() ChannelEventProcessor {
   return func(l *log.Logger, wid byte, cid byte, e interface{}) {
      if event, ok := e.(*nxPickedUpEvent); ok {
         as := getSessionForCharacterId(event.CharacterId)
         if as == nil {
            l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
            return
         }
         (*as).Announce(writer.WriteHint(fmt.Sprintf(nxGainFormat, event.Gain), 300, 10))
         (*as).Announce(writer.WriteEnableActions())
      } else {
         l.Printf("[ERROR] unable to cast event provided to handler [HandleNXPickedUpEvent]")
      }
   }
}
