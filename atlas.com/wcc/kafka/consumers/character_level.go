package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterLevelEvent struct {
	CharacterId uint32 `json:"characterId"`
}

func CharacterLevelEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterLevelEvent{}
	}
}

func HandleCharacterLevelEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterLevelEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForEachOtherSessionInMap(wid, cid, event.CharacterId, showForeignEffect(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showForeignEffect(_ logrus.FieldLogger, event *characterLevelEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteShowForeignEffect(event.CharacterId, 0))
	}
}
