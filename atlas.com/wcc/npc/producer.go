package npc

import (
	"atlas-wcc/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type setReturnTextCommand struct {
	CharacterId uint32 `json:"characterId"`
	Text        string `json:"text"`
}

type continueConversationCommand struct {
	CharacterId uint32 `json:"characterId"`
	Mode        byte   `json:"mode"`
	Type        byte   `json:"type"`
	Selection   int32  `json:"selection"`
}

type startConversationCommand struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	MapId       uint32 `json:"mapId"`
	CharacterId uint32 `json:"characterId"`
	NPCId       uint32 `json:"npcId"`
	NPCObjectId uint32 `json:"npcObjectId"`
}

func SetReturnText(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, returnText string) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_SET_RETURN_TEXT")
	return func(characterId uint32, returnText string) {
		e := &setReturnTextCommand{
			CharacterId: characterId,
			Text:        returnText,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

func ContinueConversation(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, action byte, messageType byte, selection int32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CONTINUE_NPC_CONVERSATION")
	return func(characterId uint32, action byte, messageType byte, selection int32) {
		e := &continueConversationCommand{
			CharacterId: characterId,
			Mode:        action,
			Type:        messageType,
			Selection:   selection,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}

func StartConversation(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, mapId uint32, characterId uint32, npcId uint32, objectId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_START_NPC_CONVERSATION")
	return func(worldId byte, channelId byte, mapId uint32, characterId uint32, npcId uint32, objectId uint32) {
		e := &startConversationCommand{
			WorldId:     worldId,
			ChannelId:   channelId,
			MapId:       mapId,
			CharacterId: characterId,
			NPCId:       npcId,
			NPCObjectId: objectId,
		}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}
