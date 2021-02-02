package consumers

import (
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type monsterEvent struct {
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
		return &monsterEvent{}
	}
}

func HandleMonsterEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*monsterEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			if event.Type == "CREATED" {
				createMonster(l, *event)
			} else if event.Type == "DESTROYED" {
				destroyMonster(l, *event)
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleMonsterEvent]")
		}
	}
}

func createMonster(l *log.Logger, event monsterEvent) {
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

func destroyMonster(l *log.Logger, event monsterEvent) {
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
