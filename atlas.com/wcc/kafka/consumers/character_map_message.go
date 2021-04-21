package consumers

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	"atlas-wcc/socket/response/writer"
	"github.com/sirupsen/logrus"
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
	return func(l logrus.FieldLogger, wid byte, cid byte, e interface{}) {
		if event, ok := e.(*characterMapMessageEvent); ok {
			if actingSession := processors.GetSessionByCharacterId(event.CharacterId); actingSession == nil {
				return
			}

			processors.ForEachSessionInMap(wid, cid, event.MapId, showChatText(l, event))
		} else {
			l.Errorf("Unable to cast event provided to handler")
		}
	}
}

func showChatText(_ logrus.FieldLogger, event *characterMapMessageEvent) processors.SessionOperator {
	return func(session mapleSession.MapleSession) {
		session.Announce(writer.WriteChatText(event.CharacterId, event.Message, event.GM, event.Show))
	}
}
