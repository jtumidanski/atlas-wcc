package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterMesoEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

func CharacterMesoEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterMesoEvent{}
	}
}

func HandleCharacterMesoEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterMesoEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, showMesoGain(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showMesoGain(_ logrus.FieldLogger, event *characterMesoEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteShowMesoGain(event.Gain, false))
		session.Announce(writer.WriteEnableActions())
	}
}
