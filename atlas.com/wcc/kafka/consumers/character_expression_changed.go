package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterExpressionChangedEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachSessionInMap(wid, cid, event.MapId, writeCharacterExpression(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeCharacterExpression(_ logrus.FieldLogger, event *characterExpressionChangedEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteCharacterExpression(event.CharacterId, event.Expression))
	}
}
