package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
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
			processors.ForEachSessionInMap(l, wid, cid, event.MapId, showChatText(event))
		} else {
			l.Printf("[ERROR] unable to cast event provided to handler [HandleEnableActionsEvent]")
		}
	}
}

func showChatText(event *characterMapMessageEvent) processors.SessionOperator {
	return func(l *log.Logger, session mapleSession.MapleSession) {
		session.Announce(writer.WriteChatText(event.CharacterId, event.Message, event.GM, event.Show))
	}
}
