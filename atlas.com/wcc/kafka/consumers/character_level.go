package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterLevelEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachOtherSessionInMap(wid, cid, event.CharacterId, showForeignEffect(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showForeignEffect(_ logrus.FieldLogger, event *characterLevelEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteShowForeignEffect(event.CharacterId, 0))
	}
}
