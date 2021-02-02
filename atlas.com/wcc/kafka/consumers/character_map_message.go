package consumers

import (
	"atlas-wcc/socket/response/writer"
	"log"
)

type characterMapMessageEvent struct {
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	Message     string `json:"message"`
	GM          bool   `json:"gm"`
	Show        bool   `json:"show"`
}

func CharacterMapMessageEventCreator() EmptyEventCreator {
	return func() interface{} {
		return &characterMapMessageEvent{}
	}
}

func HandleCharacterMapMessageEvent() ChannelEventProcessor {
	return func(l *log.Logger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterMapMessageEvent); ok {
			as := getSessionForCharacterId(event.CharacterId)
			if as == nil {
				l.Printf("[ERROR] unable to locate session for character %d.", event.CharacterId)
				return
			}

			sl, err := getSessionsForThoseInMap((*as).WorldId(), (*as).ChannelId(), event.MapId)
			if err != nil {
				return
			}
			for _, s := range sl {
				s.Announce(writer.WriteChatText(event.CharacterId, event.Message, event.GM, event.Show))
			}
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleEnableActionsEvent]")
		}
	}
}