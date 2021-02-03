package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
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
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*MonsterKilledEvent); ok {
			as := getSessionForCharacterId(event.KillerId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.KillerId)
				return
			}
			sl, err := getSessionsForThoseInMap(wid, cid, event.MapId)
			if err != nil {
				return
			}
			for _, s := range sl {
				s.Announce(writer.WriteKillMonster(event.UniqueId, true))
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleMonsterKilledEvent]")
		}
	}
}
