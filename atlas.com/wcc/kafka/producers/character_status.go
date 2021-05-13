package producers

import (
	"github.com/sirupsen/logrus"
)

type characterStatusEvent struct {
	WorldId     byte   `json:"worldId"`
	ChannelId   byte   `json:"channelId"`
	AccountId   uint32 `json:"accountId"`
	CharacterId uint32 `json:"characterId"`
	Type        string `json:"type"`
}

func Login(l logrus.FieldLogger) func(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_STATUS")
	return func(worldId byte, channelId byte, accountId uint32, characterId uint32) {
		emitStatus(producer, worldId, channelId, accountId, characterId, "LOGIN")
	}
}

func Logout(l logrus.FieldLogger) func(worldId byte, channelId byte, accountId uint32, characterId uint32) {
	producer := ProduceEvent(l, "TOPIC_CHARACTER_STATUS")
	return func(worldId byte, channelId byte, accountId uint32, characterId uint32) {
		emitStatus(producer, worldId, channelId, accountId, characterId, "LOGOUT")
	}
}

func emitStatus(producer func(key []byte, event interface{}), worldId byte, channelId byte, accountId uint32, characterId uint32, theType string) {
	e := &characterStatusEvent{
		WorldId:     worldId,
		ChannelId:   channelId,
		AccountId:   accountId,
		CharacterId: characterId,
		Type:        theType,
	}
	producer(CreateKey(int(characterId)), e)
}
