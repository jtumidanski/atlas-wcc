package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type DamageEntry struct {
	CharacterId uint32 `json:"characterId"`
	Damage      uint64 `json:"damage"`
}

type MonsterKilledEvent struct {
	WorldId       byte          `json:"worldId"`
	ChannelId     byte          `json:"channelId"`
	MapId         uint32        `json:"mapId"`
	UniqueId      uint32        `json:"uniqueId"`
	MonsterId     uint32        `json:"monsterId"`
	X             int16         `json:"x"`
	Y             int16         `json:"y"`
	KillerId      uint32        `json:"killerId"`
	DamageEntries []DamageEntry `json:"damageEntries"`
}

func MonsterKilledEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &MonsterKilledEvent{}
	}
}

func HandleMonsterKilledEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*MonsterKilledEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			l.Infof("Character %d killed %d.", event.UniqueId, event.KillerId)
			processors.ForEachSessionInMap(wid, cid, event.MapId, killMonster(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func killMonster(_ logrus.FieldLogger, event *MonsterKilledEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteKillMonster(event.UniqueId, true))
	}
}
