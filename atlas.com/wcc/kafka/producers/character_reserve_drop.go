package producers

import (
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

type characterReserveDropEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

func CharacterReserveDrop(l logrus.FieldLogger, span opentracing.Span) func(characterId uint32, dropId uint32) {
	producer := ProduceEvent(l, span, "TOPIC_RESERVE_DROP_COMMAND")
	return func(characterId uint32, dropId uint32) {
		e := &characterReserveDropEvent{
			CharacterId: characterId,
			DropId:      dropId,
		}
		producer(CreateKey(int(dropId)), e)
	}
}