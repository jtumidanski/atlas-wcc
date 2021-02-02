package consumers

import (
   "atlas-wcc/domain"
   "atlas-wcc/processors"
   "atlas-wcc/socket/response/writer"
   "log"
)

type MonsterEvent struct {
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	MapId     uint32 `json:"mapId"`
	UniqueId  uint32 `json:"uniqueId"`
	MonsterId uint32 `json:"monsterId"`
	ActorId   uint32 `json:"actorId"`
	Type      string `json:"type"`
}

type MonsterHandler struct {
   worldId byte
   channelId byte
}

func NewMonsterHandler(worldId byte, channelId byte) MonsterHandler {
   return MonsterHandler{worldId, channelId}
}

func (h MonsterHandler) topicToken() string {
   return "TOPIC_MONSTER_EVENT"
}

func (h MonsterHandler) emptyEventCreator() interface{} {
   return &MonsterEvent{}
}

func (h MonsterHandler) eventProcessor(l *log.Logger, event interface{}) {
   h.processEvent(l, *event.(*MonsterEvent))
}

func (h MonsterHandler) processEvent(l *log.Logger, event MonsterEvent) {
   if h.WorldId() != event.WorldId || h.ChannelId() != event.ChannelId {
      return
   }

   if event.Type == "CREATED" {
      h.created(l, event)
   } else if event.Type == "DESTROYED" {
      h.destroyed(l, event)
   }
}

func (h MonsterHandler) created(l *log.Logger, event MonsterEvent) {
   m, err := processors.GetMonster(event.UniqueId)
   if err != nil {
      return
   }
   h.create(l, *m, event)
}

func (h MonsterHandler) create(l *log.Logger, m domain.Monster, event MonsterEvent) {
   sl, err := getSessionsForThoseInMap(event.WorldId, event.ChannelId, event.MapId)
   if err != nil {
      return
   }
   for _, s := range sl {
      l.Printf("[INFO] spawning monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
      s.Announce(writer.WriteSpawnMonster(m, false))
   }
}

func (h MonsterHandler) destroyed(_ *log.Logger, event MonsterEvent) {
   sl, err := getSessionsForThoseInMap(event.WorldId, event.ChannelId, event.MapId)
   if err != nil {
      return
   }
   for _, s := range sl {
      s.Announce(writer.WriteKillMonster(event.UniqueId, false))
   }
   for _, s := range sl {
      s.Announce(writer.WriteKillMonster(event.UniqueId, true))
   }
}

func (h MonsterHandler) WorldId() byte {
   return h.worldId
}

func (h MonsterHandler) ChannelId() byte {
   return h.channelId
}
