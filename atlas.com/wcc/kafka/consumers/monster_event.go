package consumers

import (
	"atlas-wcc/domain"
	"atlas-wcc/mapleSession"
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

			monster, err := processors.GetMonster(event.UniqueId)
			if err != nil {
				l.Printf("[ERROR] unable to monster %d to create.", event.UniqueId)
				return
			}

			var handler processors.SessionOperator
			if event.Type == "CREATED" {
				handler = createMonster(event, *monster)
			} else if event.Type == "DESTROYED" {
				handler = destroyMonster(event)
			} else {
				l.Printf("[WARN] unable to handle %s event type for monster events.", event.Type)
				return
			}

			processors.ForEachSessionInMap(l, wid, cid, event.MapId, handler)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleMonsterEvent]")
		}
	}
}

func destroyMonster(event *monsterEvent) processors.SessionOperator {
	return func(l *log.Logger, session mapleSession.MapleSession) {
		session.Announce(writer.WriteKillMonster(event.UniqueId, false))
		session.Announce(writer.WriteKillMonster(event.UniqueId, true))
	}
}

func createMonster(_ *monsterEvent, monster domain.Monster) processors.SessionOperator {
	return func(l *log.Logger, session mapleSession.MapleSession) {
		session.Announce(writer.WriteSpawnMonster(monster, false))
	}
}
