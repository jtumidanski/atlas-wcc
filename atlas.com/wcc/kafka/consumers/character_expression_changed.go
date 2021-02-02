package consumers

import (
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
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}

			sl, err := getSessionsForThoseInMap((*as).WorldId(), (*as).ChannelId(), event.MapId)
			if err != nil {
				return
			}
			for _, s := range sl {
				s.Announce(writer.WriteCharacterExpression(event.CharacterId, event.Expression))
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterExpressionChangedEvent]")
		}
	}
}
