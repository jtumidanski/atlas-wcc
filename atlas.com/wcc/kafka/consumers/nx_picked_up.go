package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"fmt"
	"github.com/sirupsen/logrus"
)

const nxGainFormat = "You have earned #e#b%d NX#k#n."

type nxPickedUpEvent struct {
	CharacterId uint32 `json:"characterId"`
	Gain        uint32 `json:"gain"`
}

func NXPickedUpEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &nxPickedUpEvent{}
	}
}

func HandleNXPickedUpEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*nxPickedUpEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForSessionByCharacterId(event.CharacterId, showNXGain(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showNXGain(_ logrus.FieldLogger, event *nxPickedUpEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteHint(fmt.Sprintf(nxGainFormat, event.Gain), 300, 10))
		s.Announce(writer.WriteEnableActions())
	}
}
