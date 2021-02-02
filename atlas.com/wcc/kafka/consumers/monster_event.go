package consumers

import (
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

func MonsterEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &MonsterEvent{}
	}
}

func HandleMonsterEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, event interface{}) {
		e := *event.(*MonsterEvent)
		if wid != e.WorldId || cid != e.ChannelId {
			return
		}

		if e.Type == "CREATED" {
			createMonster(l, e)
		} else if e.Type == "DESTROYED" {
			destroyMonster(l, e)
		}
	}
}

func createMonster(l *log.Logger, event MonsterEvent) {
	m, err := processors.GetMonster(event.UniqueId)
	if err != nil {
		l.Printf("[ERROR] unable to monster %d to create.", event.UniqueId)
		return
	}

	sl, err := getSessionsForThoseInMap(event.WorldId, event.ChannelId, event.MapId)
	if err != nil {
		l.Printf("[ERROR] unable to locate sessions for map %d-%d-%d.", event.WorldId, event.ChannelId, event.MapId)
		return
	}
	for _, s := range sl {
		l.Printf("[INFO] spawning monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
		s.Announce(writer.WriteSpawnMonster(*m, false))
	}
}

func destroyMonster(l *log.Logger, event MonsterEvent) {
	sl, err := getSessionsForThoseInMap(event.WorldId, event.ChannelId, event.MapId)
	if err != nil {
		l.Printf("[ERROR] unable to locate sessions for map %d-%d-%d.", event.WorldId, event.ChannelId, event.MapId)
		return
	}
	for _, s := range sl {
		s.Announce(writer.WriteKillMonster(event.UniqueId, false))
	}
	for _, s := range sl {
		s.Announce(writer.WriteKillMonster(event.UniqueId, true))
	}
}
