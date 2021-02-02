package consumers

import (
   "atlas-wcc/socket/response/writer"
   "log"
)

type characterMesoEvent struct {
   CharacterId uint32 `json:"characterId"`
   Gain        uint32 `json:"gain"`
}

func CharacterMesoEventCreator() EmptyEventCreator {
   return func() interface{} {
      return &characterMesoEvent{}
   }
}

func HandleCharacterMesoEvent() ChannelEventProcessor {
   return func(l *log.Logger, wid byte, cid byte, e interface{}) {
      if event, ok := e.(*characterMesoEvent); ok {
         as := getSessionForCharacterId(event.CharacterId)
         if as == nil {
            l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
            return
         }
         (*as).Announce(writer.WriteShowMesoGain(event.Gain, false))
         (*as).Announce(writer.WriteEnableActions())
      } else {
         l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterMesoEvent]")
      }
   }
}
