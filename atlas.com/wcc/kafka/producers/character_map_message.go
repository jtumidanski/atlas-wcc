package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type characterMapMessageEvent struct {
	CharacterId uint32 `json:"characterId"`
	MapId       uint32 `json:"mapId"`
	Message     string `json:"message"`
	GM          bool   `json:"gm"`
	Show        bool   `json:"show"`
}

var CharacterMapMessage = func(l logrus.FieldLogger, ctx context.Context) *characterMapMessage {
	return &characterMapMessage{
		l:   l,
		ctx: ctx,
	}
}

type characterMapMessage struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func (m *characterMapMessage) Emit(worldId byte, channelId byte, mapId uint32, characterId uint32, message string, gm bool, show bool) {
	e := &characterMapMessageEvent{
		CharacterId: characterId,
		MapId:       mapId,
		Message:     message,
		GM:          gm,
		Show:        show,
	}
	produceEvent(m.l, "TOPIC_CHARACTER_MAP_MESSAGE_EVENT", createKey(int(worldId)*1000+int(channelId)), e)
}
