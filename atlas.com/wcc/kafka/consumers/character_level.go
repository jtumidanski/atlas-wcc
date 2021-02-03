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
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachOtherSessionInMap(wid, cid, event.CharacterId, showForeignEffect(l, event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterLevelEvent]")
		}
	}
}

func showForeignEffect(_ *log.Logger, event *characterLevelEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteShowForeignEffect(event.CharacterId, 0))
	}
}
