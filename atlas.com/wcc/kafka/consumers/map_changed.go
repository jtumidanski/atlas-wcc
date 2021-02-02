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

type MapChangedHandler struct {
}

func (h MapChangedHandler) topicToken() string {
   return "TOPIC_CHANGE_MAP_EVENT"
}

func (h MapChangedHandler) emptyEventCreator() interface{} {
   return &MapChangedEvent{}
}

func (h MapChangedHandler) eventProcessor(l *log.Logger, event interface{}) {
   h.processEvent(l, *event.(*MapChangedEvent))
}

func (h MapChangedHandler) processEvent(l *log.Logger, event MapChangedEvent) {
   l.Printf("[INFO] processing MapChangedEvent for character %d.", event.CharacterId)
   as := getSessionForCharacterId(event.CharacterId)
   if as == nil {
      l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
      return
   }

   catt, err := requests.GetCharacterAttributesById(event.CharacterId)
   if err != nil {
      l.Printf("[ERROR] unable to retrieve character attributes for character %d.", event.CharacterId)
      return
   }

   (*as).Announce(writer.WriteWarpToMap(event.ChannelId, event.MapId, event.PortalId, catt.Data().Attributes.Hp))
}
