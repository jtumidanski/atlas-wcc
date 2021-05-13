package consumers

import (
	"atlas-wcc/domain"
	"atlas-wcc/kafka/handler"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
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

func MonsterEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &monsterEvent{}
	}
}

func HandleMonsterEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*monsterEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			monster, err := processors.GetMonster(event.UniqueId)
			if err != nil {
				l.WithError(err).Errorf("Unable to monster %d to create.", event.UniqueId)
				return
			}

			var handler processors.SessionOperator
			if event.Type == "CREATED" {
				handler = createMonster(l, event, *monster)
			} else if event.Type == "DESTROYED" {
				handler = destroyMonster(l, event)
			} else {
				l.Warnf("Unable to handle %s event type for monster events.", event.Type)
				return
			}

			processors.ForEachSessionInMap(wid, cid, event.MapId, handler)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func destroyMonster(_ logrus.FieldLogger, event *monsterEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteKillMonster(event.UniqueId, false))
		session.Announce(writer.WriteKillMonster(event.UniqueId, true))
	}
}

func createMonster(_ logrus.FieldLogger, _ *monsterEvent, monster domain.Monster) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteSpawnMonster(monster, false))
	}
}
