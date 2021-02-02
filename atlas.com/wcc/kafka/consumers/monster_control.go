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
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] cannot location session for character %d for monster control event processing.", event.CharacterId)
				return
			}

			if event.Type == "START" {
				startControl(l, *as, *event)
			} else if event.Type == "STOP" {
				stopControl(l, *as, *event)
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleEnableActionsEvent]")
		}
	}
}

func startControl(l *log.Logger, s mapleSession.MapleSession, event monsterControlEvent) {
	m, err := processors.GetMonster(event.UniqueId)
	if err != nil {
		l.Printf("[ERROR] cannot locate monster %d for monster control event processing.", event.UniqueId)
		return
	}
	l.Printf("[INFO] controlling monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
	s.Announce(writer.WriteControlMonster(m, false, false))
}

func stopControl(l *log.Logger, s mapleSession.MapleSession, event monsterControlEvent) {
	m, err := processors.GetMonster(event.UniqueId)
	if err != nil {
		l.Printf("[ERROR] cannot locate monster %d for monster control event processing.", event.UniqueId)
		return
	}
	l.Printf("[INFO] removing control of monster %d type %d for character %d", m.UniqueId(), m.MonsterId(), s.CharacterId())
	s.Announce(writer.WriteStopControlMonster(m))
}
