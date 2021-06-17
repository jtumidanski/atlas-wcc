package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterExpressionChangedEvent struct {
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	Expression  uint32 `json:"expression"`
}

func CharacterExpressionChangedEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterExpressionChangedEvent{}
	}
}

func HandleCharacterExpressionChangedEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterExpressionChangedEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForEachSessionInMap(wid, cid, event.MapId, writeCharacterExpression(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func writeCharacterExpression(_ logrus.FieldLogger, event *characterExpressionChangedEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteCharacterExpression(event.CharacterId, event.Expression))
	}
}
