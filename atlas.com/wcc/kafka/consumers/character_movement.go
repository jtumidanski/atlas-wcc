package consumers

import (
	"atlas-wcc/rest/requests"
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
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}

			catt, err := requests.Character().GetCharacterAttributesById(event.CharacterId)
			if err != nil {
				l.Printf("[ERROR] unable to retrieve character attributes for character %d.", event.CharacterId)
				return
			}

			sl, err := getSessionsForThoseInMap(event.WorldId, event.ChannelId, catt.Data().Attributes.MapId)
			if err != nil {
				return
			}
			for _, s := range sl {
				if s.CharacterId() != event.CharacterId {
					s.Announce(writer.WriteMoveCharacter(event.CharacterId, event.RawMovement))
				}
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleEnableActionsEvent]")
		}
	}
}
