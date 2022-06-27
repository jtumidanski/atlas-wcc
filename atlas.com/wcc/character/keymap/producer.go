package keymap

import (
	"atlas-wcc/kafka"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type changeKeyMapCommand struct {
	CharacterId uint32   `json:"characterId"`
	Changes     []Change `json:"changes"`
}

type Change struct {
	Key        int32 `json:"key"`
	ChangeType int8  `json:"changeType"`
	Action     int32 `json:"action"`
}

func ChangeKeyMap(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, changes []Change) {
	producer := kafka.ProduceEvent(l, span, "TOPIC_CHANGE_KEY_MAP")
	return func(characterId uint32, changes []Change) {
		e := changeKeyMapCommand{CharacterId: characterId, Changes: changes}
		producer(kafka.CreateKey(int(characterId)), e)
	}
}
