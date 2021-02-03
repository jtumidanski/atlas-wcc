package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type monsterControlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	UniqueId    uint32 `json:"uniqueId"`
	Type        string `json:"type"`
}

func MonsterControlEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &monsterControlEvent{}
	}
}

func HandleMonsterControlEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*monsterControlEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			var handler processors.SessionOperator
			if event.Type == "START" {
				handler = startControl(l, event)
			} else if event.Type == "STOP" {
				handler = stopControl(l, event)
			} else {
				l.Printf("[WARN] received unhandled monster control event type of %s", event.Type)
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, handler)
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleEnableActionsEvent]")
		}
	}
}

func stopControl(_ *log.Logger, event *monsterControlEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		m, err := processors.GetMonster(event.UniqueId)
		if err != nil {
			return
		}
		session.Announce(writer.WriteStopControlMonster(m))
	}
}

func startControl(_ *log.Logger, event *monsterControlEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		m, err := processors.GetMonster(event.UniqueId)
		if err != nil {
			return
		}
		session.Announce(writer.WriteControlMonster(m, false, false))
	}
}
