package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
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
			processors.ForSessionByCharacterId(l, event.CharacterId, warpCharacter(event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleChangeMapEvent]")
		}
	}
}

func warpCharacter(event *mapChangedEvent) processors.SessionOperator {
	return func(l *log.Logger, session mapleSession.MapleSession) {
		catt, err := requests.Character().GetCharacterAttributesById(event.CharacterId)
		if err != nil {
			return
		}
		session.Announce(writer.WriteWarpToMap(event.ChannelId, event.MapId, event.PortalId, catt.Data().Attributes.Hp))
	}
}
