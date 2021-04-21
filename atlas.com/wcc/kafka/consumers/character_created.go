package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
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

func CharacterCreatedEventCreator() EmptyEventCreator {
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

			processors.ForEachGMSession(announceCharacterCreated(event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func announceCharacterCreated(event *characterCreatedEvent) processors.SessionOperator {
	return func(sessions mapleSession.MapleSession) {
		sessions.Announce(writer.WriteYellowTip(fmt.Sprintf(characterCreatedFormat, event.Name)))
	}
}
