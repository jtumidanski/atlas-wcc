package consumers

import (
   "atlas-wcc/mapleSession"
   "atlas-wcc/processors"
   "atlas-wcc/socket/response/writer"
   "log"
)

type MonsterControlEvent struct {
   WorldId     byte   `json:"worldId"`
   ChannelId   byte   `json:"channelId"`
   CharacterId uint32 `json:"characterId"`
   UniqueId    uint32 `json:"uniqueId"`
   Type        string `json:"type"`
}

type MonsterControlHandler struct {
}

func (h MonsterControlHandler) topicToken() string {
   return "TOPIC_CONTROL_MONSTER_EVENT"
}

func (h MonsterControlHandler) emptyEventCreator() interface{} {
   return &MonsterControlEvent{}
}

func (h MonsterControlHandler) eventProcessor(l *log.Logger, event interface{}) {
   h.processEvent(l, *event.(*MonsterControlEvent))
}

func (h MonsterControlHandler) processEvent(l *log.Logger, event MonsterControlEvent) {
   as := getSessionForCharacterId(event.CharacterId)
   if as == nil {
      l.Printf("[ERROR] cannot location session for character %d for monster control event processing.", event.CharacterId)
      return
   }

   if event.Type == "START" {
      h.start(l, *as, event)
   } else if event.Type == "STOP" {
      h.stop(l, *as, event)
   }
}

func (h MonsterControlHandler) start(l *log.Logger, s mapleSession.MapleSession, event MonsterControlEvent) {
   m, err := processors.GetMonster(event.UniqueId)
   if err != nil {
      l.Printf("[ERROR] cannot locate monster %d for monster control event processing.", event.UniqueId)
      return
   }
   l.Printf("[INFO] controlling monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
   s.Announce(writer.WriteControlMonster(m, false, false))
}

func (h MonsterControlHandler) stop(l *log.Logger, s mapleSession.MapleSession, event MonsterControlEvent) {
   m, err := processors.GetMonster(event.UniqueId)
   if err != nil {
      l.Printf("[ERROR] cannot locate monster %d for monster control event processing.", event.UniqueId)
      return
   }
   l.Printf("[INFO] removing control of monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
   s.Announce(writer.WriteStopControlMonster(m))
}
