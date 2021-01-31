package handler

import (
   "atlas-wcc/kafka/producers"
   "atlas-wcc/mapleSession"
   "atlas-wcc/processors"
   "context"
   "github.com/jtumidanski/atlas-socket/request"
   "log"
)

const OpChangeMap uint16 = 0x26

type ChangeMapRequest struct {
   cashShop  bool
   fromDying int8
   targetId  int32
   startWarp string
   wheel     bool
}

func (r ChangeMapRequest) CashShop() bool {
   return r.cashShop
}

func (r ChangeMapRequest) StartWarp() string {
   return r.startWarp
}

func ReadChangeMapRequest(reader *request.RequestReader) ChangeMapRequest {
   cs := len(reader.String()) == 0
   fromDying := int8(-1)
   targetId := int32(-1)
   startWarp := ""
   wheel := false

   if !cs {
      fromDying = reader.ReadInt8()
      targetId = reader.ReadInt32()
      startWarp = reader.ReadAsciiString()
      reader.ReadByte()
      wheel = reader.ReadInt16() > 0
   }
   return ChangeMapRequest{cs, fromDying, targetId, startWarp, wheel}
}

type ChangeMapHandler struct {
}

func (h *ChangeMapHandler) IsValid(l *log.Logger, ms *mapleSession.MapleSession) bool {
   v := processors.IsLoggedIn((*ms).AccountId())
   if !v {
      l.Printf("[ERROR] attempting to process a [ChangeMapRequest] when the account %d is not logged in.", (*ms).SessionId())
   }
   return v
}

func (h *ChangeMapHandler) HandleRequest(l *log.Logger, s *mapleSession.MapleSession, r *request.RequestReader) {
   p := ReadChangeMapRequest(r)
   if p.CashShop() {

   } else {
      ca, err := processors.GetCharacterAttributesById((*s).CharacterId())
      if err != nil {
         return
      }

      portal, err := processors.GetPortalByName(ca.MapId(), p.StartWarp())
      if err != nil {
         return
      }
      producers.NewPortalEnter(l, context.Background()).EmitEnter((*s).WorldId(), (*s).ChannelId(), ca.MapId(), portal.Id(), (*s).CharacterId())
   }
}
