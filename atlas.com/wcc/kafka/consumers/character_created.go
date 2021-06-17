package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"fmt"
	"github.com/sirupsen/logrus"
)

const characterCreatedFormat = "Character %s has been created."

type characterCreatedEvent struct {
	WorldId     byte   `json:"worldId"`
	CharacterId uint32 `json:"characterId"`
	Name        string `json:"name"`
}

func CharacterCreatedEventCreator() handler.EmptyEventCreator {
	return func() interface{} {
		return &characterCreatedEvent{}
	}
}

func HandleCharacterCreatedEvent() ChannelEventProcessor {
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterCreatedEvent); ok {
			if wid != event.WorldId {
				return
			}

			session.ForEachGM(announceCharacterCreated(l)(event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func announceCharacterCreated(l logrus.FieldLogger) func(event *characterCreatedEvent) session.Operator {
	return func(event *characterCreatedEvent) session.Operator {
		b := writer.WriteYellowTip(fmt.Sprintf(characterCreatedFormat, event.Name))
		return func(s *session.Model) {
			err := s.Announce(b)
			if err != nil {
				l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
			}
		}
	}
}
