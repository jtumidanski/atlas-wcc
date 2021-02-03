package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
)

type characterExpressionChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	Expression  uint32 `json:"expression"`
}

func CharacterExpressionChangedEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterExpressionChangedEvent{}
	}
}

func HandleCharacterExpressionChangedEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterExpressionChangedEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachSessionInMap(wid, cid, event.MapId, writeCharacterExpression(l, event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterExpressionChangedEvent]")
		}
	}
}

func writeCharacterExpression(_ *log.Logger, event *characterExpressionChangedEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteCharacterExpression(event.CharacterId, event.Expression))
	}
}
