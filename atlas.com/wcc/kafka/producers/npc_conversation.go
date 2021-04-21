package producers

import (
	"context"
	"github.com/sirupsen/logrus"
)

type setReturnTextCommand struct {
	CharacterId uint32 `json:"characterId"`
	Text        string `json:"text"`
}

type continueNPCConversationCommand struct {
	CharacterId uint32 `json:"characterId"`
	Mode        byte   `json:"mode"`
	Type        byte   `json:"type"`
	Selection   int32  `json:"selection"`
}

type startNPCConversationCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	NPCId       uint32 `json:"npcId"`
	NPCObjectId uint32 `json:"npcObjectId"`
}

var NPCConversation = func(l logrus.FieldLogger, ctx context.Context) *npcConversation {
	return &npcConversation{
		l:   l,
		ctx: ctx,
	}
}

type npcConversation struct {
	l   logrus.FieldLogger
	ctx context.Context
}

func (m *npcConversation) SetReturnText(characterId uint32, returnText string) {
	e := &setReturnTextCommand{
		CharacterId: characterId,
		Text:        returnText,
	}
	produceEvent(m.l, "TOPIC_SET_RETURN_TEXT", createKey(int(characterId)), e)
}

func (m *npcConversation) ContinueConversation(characterId uint32, action byte, messageType byte, selection int32) {
	e := &continueNPCConversationCommand{
		CharacterId: characterId,
		Mode:        action,
		Type:        messageType,
		Selection:   selection,
	}
	produceEvent(m.l, "TOPIC_CONTINUE_NPC_CONVERSATION", createKey(int(characterId)), e)
}

func (m *npcConversation) StartConversation(worldId byte, channelId byte, mapId uint32, characterId uint32, npcId uint32, objectId uint32) {
	e := &startNPCConversationCommand{
		WorldId:     worldId,
		ChannelId:   channelId,
		MapId:       mapId,
		CharacterId: characterId,
		NPCId:       npcId,
		NPCObjectId: objectId,
	}
	produceEvent(m.l, "TOPIC_START_NPC_CONVERSATION", createKey(int(characterId)), e)
}
