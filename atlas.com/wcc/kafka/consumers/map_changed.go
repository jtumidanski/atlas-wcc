package consumers

import (
	"atlas-wcc/kafka/handler"
	"atlas-wcc/rest/requests"
	"atlas-wcc/session"
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

func ChangeMapEventCreator() handler.EmptyEventCreator {
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

			session.ForSessionByCharacterId(event.CharacterId, warpCharacter(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func warpCharacter(l logrus.FieldLogger, event *mapChangedEvent) session.Operator {
	return func(s *session.Model) {
		catt, err := requests.Character().GetCharacterAttributesById(event.CharacterId)
		if err != nil {
			l.WithError(err).Errorf("Unable to retrieve character %d properties", event.CharacterId)
			return
		}
		err = s.Announce(writer.WriteWarpToMap(l)(event.ChannelId, event.MapId, event.PortalId, catt.Data().Attributes.Hp))
		if err != nil {
			l.WithError(err).Errorf("Unable to warp character %d to map %d", event.CharacterId, event.MapId)
		}
	}
}
