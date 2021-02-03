package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"fmt"
	"log"
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
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterCreatedEvent); ok {
			processors.ForEachGMSession(announceCharacterCreated(event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterCreatedEvent]")
		}
	}
}

func announceCharacterCreated(event *characterCreatedEvent) processors.SessionOperator {
	return func(sessions mapleSession.MapleSession) {
		sessions.Announce(writer.WriteYellowTip(fmt.Sprintf(characterCreatedFormat, event.Name)))
	}
}
