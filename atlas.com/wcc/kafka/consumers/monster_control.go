package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type monsterControlEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
	UniqueId    uint32 `json:"uniqueId"`
	Type        string `json:"type"`
}

func MonsterControlEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &monsterControlEvent{}
	}
}

func HandleMonsterControlEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
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
				l.Warnf("Received unhandled monster control event type of %s", event.Type)
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, handler)
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func stopControl(l logrus.FieldLogger, event *monsterControlEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		m, err := processors.GetMonster(event.UniqueId)
		if err != nil {
			return
		}
		l.Infof("Stopping control of %d for character %d.", event.UniqueId, event.CharacterId)
		session.Announce(writer.WriteStopControlMonster(m))
	}
}

func startControl(l logrus.FieldLogger, event *monsterControlEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		m, err := processors.GetMonster(event.UniqueId)
		if err != nil {
			return
		}
		l.Infof("Starting control of %d for character %d.", event.UniqueId, event.CharacterId)
		session.Announce(writer.WriteControlMonster(m, false, false))
	}
}
