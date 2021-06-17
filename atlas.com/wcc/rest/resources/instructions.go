package resources

import (
	"atlas-wcc/json"
	"atlas-wcc/rest/attributes"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"strconv"
)

type InstructionResource struct {
	l logrus.FieldLogger
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func NewInstructionResource(l logrus.FieldLogger) *InstructionResource {
	return &InstructionResource{l}
}

func (i *InstructionResource) CreateInstruction(rw http.ResponseWriter, r *http.Request) {
	characterId := getCharacterId(r)

	cs := &attributes.InstructionInputDataContainer{}
	err := json.FromJSON(cs, r.Body)
	if err != nil {
		i.l.WithError(err).Errorf("Deserializing instruction")
		rw.WriteHeader(http.StatusBadRequest)
		json.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	s := session.GetSessionByCharacterId(characterId)
	if s == nil {
		i.l.WithError(err).Errorf("Cannot locate session for instruction")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	(*s).Announce(writer.WriteHint(cs.Data.Attributes.Message, cs.Data.Attributes.Width, cs.Data.Attributes.Height))
	(*s).Announce(writer.WriteEnableActions())

	rw.WriteHeader(http.StatusNoContent)
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
