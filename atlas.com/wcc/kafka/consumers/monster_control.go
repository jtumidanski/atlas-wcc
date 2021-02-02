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

func MonsterControlEventCreator() EmptyEventCreator {
	return func() interface{} {
		return MonsterControlEvent{}
	}
}

func HandleMonsterControlEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, event interface{}) {
		e := *event.(*MonsterControlEvent)
		as := getSessionForCharacterId(e.CharacterId)
		if as == nil {
			l.Printf("[ERROR] cannot location session for character %d for monster control event processing.", e.CharacterId)
			return
		}

		if e.Type == "START" {
			startControl(l, *as, e)
		} else if e.Type == "STOP" {
			stopControl(l, *as, e)
		}
	}
}

func startControl(l *log.Logger, s mapleSession.MapleSession, event MonsterControlEvent) {
	m, err := processors.GetMonster(event.UniqueId)
	if err != nil {
		l.Printf("[ERROR] cannot locate monster %d for monster control event processing.", event.UniqueId)
		return
	}
	l.Printf("[INFO] controlling monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
	s.Announce(writer.WriteControlMonster(m, false, false))
}

func stopControl(l *log.Logger, s mapleSession.MapleSession, event MonsterControlEvent) {
	m, err := processors.GetMonster(event.UniqueId)
	if err != nil {
		l.Printf("[ERROR] cannot locate monster %d for monster control event processing.", event.UniqueId)
		return
	}
	l.Printf("[INFO] removing control of monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
	s.Announce(writer.WriteStopControlMonster(m))
}
