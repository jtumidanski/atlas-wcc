package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/rest/requests"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*mapChangedEvent); ok {
			if wid != event.WorldId || cid != event.ChannelId {
				return
			}

			processors.ForSessionByCharacterId(event.CharacterId, warpCharacter(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func warpCharacter(_ logrus.FieldLogger, event *mapChangedEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		catt, err := requests.Character().GetCharacterAttributesById(event.CharacterId)
		if err != nil {
			return
		}
		session.Announce(writer.WriteWarpToMap(event.ChannelId, event.MapId, event.PortalId, catt.Data().Attributes.Hp))
	}
}
