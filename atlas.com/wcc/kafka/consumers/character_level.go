package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type characterLevelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func CharacterLevelEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterLevelEvent{}
	}
}

func HandleCharacterLevelEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterLevelEvent); ok {
			processors.ForEachOtherSessionInMap(l, wid, cid, event.CharacterId, showForeignEffect(event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterLevelEvent]")
		}
	}
}

func showForeignEffect(event *characterLevelEvent) processors.SessionOperator {
	return func(l *log.Logger, session mapleSession.MapleSession) {
		session.Announce(writer.WriteShowForeignEffect(event.CharacterId, 0))
	}
}
