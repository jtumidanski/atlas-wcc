package handler

import (
   "atlas-wcc/kafka/producers"
   "atlas-wcc/mapleSession"
   "atlas-wcc/processors"
   "context"
   "github.com/jtumidanski/atlas-socket/request"
   "log"
)

const OpChangeMapSpecial uint16 = 0x64

type ChangeMapSpecialRequest struct {
   startWarp string
}

func (c *ChangeMapSpecialRequest) StartWarp() string {
   return c.startWarp
}

func ReadChangeMapSpecialRequest(reader *request.RequestReader) ChangeMapSpecialRequest {
   reader.ReadByte()
   sw := reader.ReadAsciiString()
   reader.ReadUint16()
   return ChangeMapSpecialRequest{sw}
}

type ChangeMapSpecialHandler struct {
}

func (h *ChangeMapSpecialHandler) IsValid(l *log.Logger, ms *mapleSession.MapleSession) bool {
   v := processors.IsLoggedIn((*ms).AccountId())
   if !v {
      l.Printf("[ERROR] attempting to process a [ChangeMapSpecialRequest] when the account %d is not logged in.", (*ms).SessionId())
   }
   return v
}

func (h *ChangeMapSpecialHandler) HandleRequest(l *log.Logger, ms *mapleSession.MapleSession, r *request.RequestReader) {
   p := ReadChangeMapSpecialRequest(r)
   c, err := processors.GetCharacterAttributesById((*ms).CharacterId())
   if err != nil {
      return
   }

   portal, err := processors.GetPortalByName(c.MapId(), p.StartWarp())
   if err != nil {
      return
   }
   producers.NewPortalEnter(l, context.Background()).EmitEnter((*ms).WorldId(), (*ms).ChannelId(), c.MapId(), portal.Id(), (*ms).CharacterId())
}
