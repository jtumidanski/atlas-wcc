package consumers

import (
   "atlas-wcc/rest/requests"
   "atlas-wcc/socket/response/writer"
   "log"
)

type MapChangedEvent struct {
   WorldId     byte   `json:"worldId"`
   ChannelId   byte   `json:"channelId"`
   MapId       uint32 `json:"mapId"`
   PortalId    uint32 `json:"portalId"`
   CharacterId uint32 `json:"characterId"`
}

func ChangeMapEventCreator() EmptyEventCreator {
   return func() interface{} {
      return &MapChangedEvent{}
   }
}

func HandleChangeMapEvent() ChannelEventProcessor {
   return func(l *log.Logger, wid byte, cid byte, event interface{}) {
      e := *event.(*MapChangedEvent)

      l.Printf("[INFO] processing MapChangedEvent for character %d.", e.CharacterId)
      as := getSessionForCharacterId(e.CharacterId)
      if as == nil {
         l.Printf("[ERROR] unable to locate session for character %d.", e.CharacterId)
         return
      }

      catt, err := requests.GetCharacterAttributesById(e.CharacterId)
      if err != nil {
         l.Printf("[ERROR] unable to retrieve character attributes for character %d.", e.CharacterId)
         return
      }

      (*as).Announce(writer.WriteWarpToMap(e.ChannelId, e.MapId, e.PortalId, catt.Data().Attributes.Hp))
   }
}
