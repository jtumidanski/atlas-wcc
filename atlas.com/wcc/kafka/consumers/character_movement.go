package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
)

type characterMovementEvent struct {
	WorldId     byte        `json:"worldId"`
	ChannelId   byte        `json:"channelId"`
	CharacterId uint32      `json:"characterId"`
	X           int16       `json:"x"`
	Y           int16       `json:"y"`
	Stance      byte        `json:"stance"`
	RawMovement RawMovement `json:"rawMovement"`
}

func CharacterMovementEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterMovementEvent{}
	}
}

func HandleCharacterMovementEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterMovementEvent); ok {
			if actingSession := session.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			session.ForEachOtherSessionInMap(wid, cid, event.CharacterId, moveCharacter(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func moveCharacter(_ logrus.FieldLogger, event *characterMovementEvent) session.SessionOperator {
	return func(s session.Model) {
		s.Announce(writer.WriteMoveCharacter(event.CharacterId, event.RawMovement))
	}
}
