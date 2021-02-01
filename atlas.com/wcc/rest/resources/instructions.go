package resources

import (
   "atlas-wcc/mapleSession"
   "atlas-wcc/registries"
   "atlas-wcc/rest/attributes"
   "atlas-wcc/socket/response/writer"
   "github.com/gorilla/mux"
   "log"
   "net/http"
   "strconv"
)

type InstructionInputDataContainer struct {
   Data InstructionData `json:"data"`
}

type InstructionData struct {
   Id         string                `json:"id"`
   Type       string                `json:"type"`
   Attributes InstructionAttributes `json:"attributes"`
}

type InstructionAttributes struct {
   Message string `json:"message"`
   Width   int16 `json:"width"`
   Height  int16 `json:"height"`
}

type InstructionResource struct {
   l *log.Logger
}

// GenericError is a generic error message returned by a server
type GenericError struct {
   Message string `json:"message"`
}

func NewInstructionResource(l *log.Logger) *InstructionResource {
   return &InstructionResource{l}
}

func (i *InstructionResource) CreateInstruction(rw http.ResponseWriter, r *http.Request) {
   characterId := getCharacterId(r)

   cs := &InstructionInputDataContainer{}
   err := attributes.FromJSON(cs, r.Body)
   if err != nil {
      i.l.Println("[ERROR] deserializing instruction", err)
      rw.WriteHeader(http.StatusBadRequest)
      attributes.ToJSON(&GenericError{Message: err.Error()}, rw)
      return
   }

   s := getSessionForCharacterId(characterId)
   if s == nil {
      i.l.Println("[ERROR] cannot locate session for instruction", err)
      rw.WriteHeader(http.StatusBadRequest)
      return
   }

   (*s).Announce(writer.WriteHint(cs.Data.Attributes.Message, cs.Data.Attributes.Width, cs.Data.Attributes.Height))
   (*s).Announce(writer.WriteEnableActions())

   rw.WriteHeader(http.StatusNoContent)
}

func getSessionForCharacterId(cid uint32) *mapleSession.MapleSession {
   for _, s := range registries.GetSessionRegistry().GetAll() {
      if cid == s.CharacterId() {
         return &s
      }
   }
   return nil
}

func getWorldId(r *http.Request) byte {
   vars := mux.Vars(r)
   value, err := strconv.Atoi(vars["worldId"])
   if err != nil {
      log.Println("Error parsing worldId as byte")
      return 0
   }
   return byte(value)
}

func getChannelId(r *http.Request) byte {
   vars := mux.Vars(r)
   value, err := strconv.Atoi(vars["channelId"])
   if err != nil {
      log.Println("Error parsing channelId as byte")
      return 0
   }
   return byte(value)
}

func getCharacterId(r *http.Request) uint32 {
   vars := mux.Vars(r)
   value, err := strconv.Atoi(vars["characterId"])
   if err != nil {
      log.Println("Error parsing characterId as uint32")
      return 0
   }
   return uint32(value)
}