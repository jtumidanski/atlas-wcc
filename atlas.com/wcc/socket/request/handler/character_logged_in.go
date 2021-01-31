package handler

import (
   "atlas-wcc/kafka/producers"
   "atlas-wcc/mapleSession"
   "atlas-wcc/processors"
   "atlas-wcc/socket/response/writer"
   "context"
   "github.com/jtumidanski/atlas-socket/request"
   "log"
)

const OpCharacterLoggedIn uint16 = 0x14

type CharacterLoggedInRequest struct {
   characterId uint32
}

func (c *CharacterLoggedInRequest) CharacterId() uint32 {
   return c.characterId
}


func ReadCharacterLoggedInRequest(reader *request.RequestReader) CharacterLoggedInRequest {
   cid := reader.ReadUint32()
   return CharacterLoggedInRequest{cid}
}

type CharacterLoggedInHandler struct {
}

func (h *CharacterLoggedInHandler) IsValid(l *log.Logger, ms *mapleSession.MapleSession) bool {
   return true
}

func (h *CharacterLoggedInHandler) HandleRequest(l *log.Logger, ms *mapleSession.MapleSession, r *request.RequestReader) {
   p := ReadCharacterLoggedInRequest(r)
   c, err := processors.GetCharacterById(p.CharacterId())
   if err != nil {
      return
   }

   (*ms).SetAccountId(c.Attributes().AccountId())
   (*ms).SetCharacterId(c.Attributes().Id())
   (*ms).SetGm(c.Attributes().Gm())

   producers.NewCharacterStatus(l, context.Background()).EmitLogin((*ms).WorldId(), (*ms).ChannelId(), (*ms).AccountId(), p.CharacterId())
   _ = (*ms).Announce(writer.WriteGetCharacterInfo((*ms).ChannelId(), *c))
}
