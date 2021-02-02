package consumers

import (
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"log"
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
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterLevelEvent); ok {
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}

			catt, err := processors.GetCharacterAttributesById(event.CharacterId)
			if err != nil {
				l.Printf("[ERROR] unable to retrieve character attributes for character %d.", event.CharacterId)
				return
			}

			sl, err := getSessionsForThoseInMap(wid, cid, catt.MapId())
			if err != nil {
				return
			}
			for _, s := range sl {
				if s.CharacterId() != event.CharacterId {
					s.Announce(writer.WriteShowForeignEffect(event.CharacterId, 0))
				}
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleCharacterLevelEvent]")
		}
	}
}
