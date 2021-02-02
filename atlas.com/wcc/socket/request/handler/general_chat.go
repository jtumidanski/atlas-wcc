package handler

import (
	"atlas-wcc/kafka/producers"
	"atlas-wcc/mapleSession"
	"atlas-wcc/processors"
	request2 "atlas-wcc/socket/request"
	"context"
	"github.com/jtumidanski/atlas-socket/request"
	"log"
)

const OpGeneralChat uint16 = 0x31

type generalChatRequest struct {
	message string
	show    byte
}

func (r generalChatRequest) Message() string {
	return r.message
}

func (r generalChatRequest) Show() byte {
	return r.show
}

func readGeneralChatRequest(reader *request.RequestReader) generalChatRequest {
	message := reader.ReadAsciiString()
	show := reader.ReadByte()
	return generalChatRequest{message, show}
}

func GeneralChatHandler() request2.SessionRequestHandler {
	return func(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader) {
		p := readGeneralChatRequest(r)
		ca, err := processors.GetCharacterAttributesById((*s).CharacterId())
		if err != nil {
			l.Printf("[ERROR] cannot handle [GeneralChatRequest] because the acting character %d cannot be located.", (*s).CharacterId())
			return
		}

		producers.CharacterMapMessage(l, context.Background()).Emit((*s).WorldId(), (*s).ChannelId(), ca.MapId(), (*s).CharacterId(), p.Message(), ca.Gm(), p.Show() == 1)
	}
}