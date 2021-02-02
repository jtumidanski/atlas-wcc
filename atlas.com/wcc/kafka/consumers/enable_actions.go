package consumers

import (
   "atlas-wcc/socket/response/writer"
   "log"
)

type EnableActionsEvent struct {
   CharacterId uint32 `json:"characterId"`
}

type EnableActionsHandler struct {
}

func (h EnableActionsHandler) topicToken() string {
   return "TOPIC_ENABLE_ACTIONS"
}

func (h EnableActionsHandler) emptyEventCreator() interface{} {
   return &EnableActionsEvent{}
}

func (h EnableActionsHandler) eventProcessor(l *log.Logger, event interface{}) {
   h.processEvent(l, *event.(*EnableActionsEvent))
}

func (h EnableActionsHandler) processEvent(_ *log.Logger, event EnableActionsEvent) {
   as := getSessionForCharacterId(event.CharacterId)
   if as == nil {
      return
   }

   (*as).Announce(writer.WriteEnableActions())
}
