package consumers

import (
   "atlas-wcc/rest/requests"
   "atlas-wcc/socket/response/writer"
   "log"
)

type MonsterMovementEvent struct {
   UniqueId      uint32      `json:"uniqueId"`
   ObserverId    uint32      `json:"observerId"`
   SkillPossible bool        `json:"skillPossible"`
   Skill         int8        `json:"skill"`
   SkillId       byte        `json:"skillId"`
   SkillLevel    byte        `json:"skillLevel"`
   Option        uint16      `json:"option"`
   StartX        int16       `json:"startX"`
   StartY        int16       `json:"startY"`
   EndX          int16       `json:"endX"`
   EndY          int16       `json:"endY"`
   Stance        byte        `json:"stance"`
   RawMovement   RawMovement `json:"rawMovement"`
}

type RawMovement []byte

type MonsterMovementHandler struct {
   worldId   byte
   channelId byte
}

func NewMonsterMovementHandler(worldId byte, channelId byte) MonsterMovementHandler {
   return MonsterMovementHandler{worldId, channelId}
}

func (h MonsterMovementHandler) topicToken() string {
   return "TOPIC_MONSTER_MOVEMENT"
}

func (h MonsterMovementHandler) emptyEventCreator() interface{} {
   return &MonsterMovementEvent{}
}

func (h MonsterMovementHandler) eventProcessor(l *log.Logger, event interface{}) {
   h.processEvent(l, *event.(*MonsterMovementEvent))
}

func (h MonsterMovementHandler) processEvent(l *log.Logger, event MonsterMovementEvent) {
   m, err := requests.GetMonster(event.UniqueId)
   if err != nil {
      l.Printf("[ERROR] unable to retrieve monster %d for MonsterMovementEvent", event.UniqueId)
      return
   }

   sl, err := getSessionsForThoseInMap(h.worldId, h.channelId, m.Data().Attributes.MapId)
   if err != nil {
      return
   }
   for _, s := range sl {
      if s.CharacterId() != event.ObserverId {
         s.Announce(writer.WriteMoveMonster(event.ObserverId, event.SkillPossible, event.Skill, event.SkillId, event.SkillLevel, event.Option, event.StartX, event.StartY, event.RawMovement))
      }
   }
}
