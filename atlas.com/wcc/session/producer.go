package session

import (
	"atlas-wcc/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type statusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func Login(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_STATUS")
	return func(worldId byte, channelId byte, accountId uint32, characterId uint32) {
		emitStatus(producer, worldId, channelId, accountId, characterId, "LOGIN")
	}
}

func Logout(l logrus.FieldLogger, span opentracing.Span) func(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHARACTER_STATUS")
	return func(worldId byte, channelId byte, accountId uint32, characterId uint32) {
		emitStatus(producer, worldId, channelId, accountId, characterId, "LOGOUT")
	}
}

func emitStatus(producer func(key []byte, event interface{}), worldId byte, channelId byte, accountId uint32, characterId uint32, theType string) {
	e := &statusEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		AccountId:   accountId,
		CharacterId: characterId,
		Type:        theType,
	}
	producer(kafka.CreateKey(int(characterId)), e)
}
