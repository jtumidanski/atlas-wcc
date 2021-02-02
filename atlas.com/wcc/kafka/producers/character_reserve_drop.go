package producers

import (
	"context"
	"log"
)

type characterReserveDropEvent struct {
	CharacterId uint32 `json:"characterId"`
	DropId      uint32 `json:"dropId"`
}

var CharacterReserveDrop = func(l *log.Logger, ctx context.Context) *characterReserveDrop {
	return &characterReserveDrop{
		l:   l,
		ctx: ctx,
	}
}

type characterReserveDrop struct {
	l   *log.Logger
	ctx context.Context
}

func (m *characterReserveDrop) Emit(characterId uint32, dropId uint32) {
	e := &characterReserveDropEvent{
		CharacterId: characterId,
		DropId:      dropId,
	}
	produceEvent(m.l, "TOPIC_RESERVE_DROP_COMMAND", createKey(int(dropId)), e)
}
