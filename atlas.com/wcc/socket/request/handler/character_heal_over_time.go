package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	request2 "atlas-wcc/socket/request"
	"context"
	"github.com/jtumidanski/atlas-socket/request"
	"log"
)

const OpCharacterHealOverTime uint16 = 0x59

type healOverTimeRequest struct {
	hp uint16
	mp uint16
}

func (r healOverTimeRequest) HP() uint16 {
	return r.hp
}

func (r healOverTimeRequest) MP() uint16 {
	return r.mp
}

func readHealOverTimeRequest(reader *request.RequestReader) healOverTimeRequest {
	reader.Skip(8)
	hp := reader.ReadUint16()
	mp := reader.ReadUint16()
	return healOverTimeRequest{hp, mp}
}

func HealOverTimeHandler() request2.SessionRequestHandler {
	return func(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readHealOverTimeRequest(r)

		producers.CharacterAdjustHealth(l, context.Background()).Emit((*s).CharacterId(), p.HP())
		producers.CharacterAdjustMana(l, context.Background()).Emit((*s).CharacterId(), p.MP())
	}
}
