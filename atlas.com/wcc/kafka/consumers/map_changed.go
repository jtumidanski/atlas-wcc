package consumers

import (
	"atlas-wcc/rest/requests"
	"atlas-wcc/socket/response/writer"
	"log"
)

type mapChangedEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	PortalId    uint32 `json:"portalId"`
	CharacterId uint32 `json:"characterId"`
}

func ChangeMapEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &mapChangedEvent{}
	}
}

func HandleChangeMapEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*mapChangedEvent); ok {
			l.Printf("[INFO] processing MapChangedEvent for character %d.", event.CharacterId)
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

			(*as).Announce(writer.WriteWarpToMap(event.ChannelId, event.MapId, event.PortalId, catt.Data().Attributes.Hp))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleChangeMapEvent]")
		}
	}
}
