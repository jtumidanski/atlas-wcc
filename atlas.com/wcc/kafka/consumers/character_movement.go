package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
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

func CharacterMovementEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterMovementEvent{}
	}
}

func HandleCharacterMovementEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterMovementEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachOtherSessionInMap(wid, cid, event.CharacterId, moveCharacter(l, event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleEnableActionsEvent]")
		}
	}
}

func moveCharacter(_ *log.Logger, event *characterMovementEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteMoveCharacter(event.CharacterId, event.RawMovement))
	}
}
