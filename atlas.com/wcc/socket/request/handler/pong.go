package handler

import (
	"atlas-wcc/mapleSession"
	request2 "atlas-wcc/socket/request"
	"github.com/jtumidanski/atlas-socket/request"
	"log"
)

const OpCodePong uint16 = 0x18

func PongHandler() request2.SessionRequestHandler {
	return func(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader) {
	}
}
