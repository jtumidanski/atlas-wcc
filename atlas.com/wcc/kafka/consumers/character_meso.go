package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterMesoEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        int32  `json:"gain"`
}

func CharacterMesoEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterMesoEvent{}
	}
}

func HandleCharacterMesoEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterMesoEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, showMesoGain(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showMesoGain(_ logrus.FieldLogger, event *characterMesoEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteShowMesoGain(event.Gain, false))
		s.Announce(writer.WriteEnableActions())
	}
}
