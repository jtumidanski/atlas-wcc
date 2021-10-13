package instruction

import (
	"atlas-wcc/json"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	iRouter := router.PathPrefix("/characters/{characterId}/instructions").Subrouter()
	iRouter.HandleFunc("", HandleCreateInstruction(l)).Methods(http.MethodPost)
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

func HandleCreateInstruction(l logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		characterId := getCharacterId(l)(r)

		cs := &InputDataContainer{}
		err := json.FromJSON(cs, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Deserializing instruction")
			w.WriteHeader(http.StatusBadRequest)
			err := json.ToJSON(&GenericError{Message: err.Error()}, w)
			if err != nil {
				l.WithError(err).Errorf("Unable to serialize error mesage")
			}
			return
		}

		s := session.GetByCharacterId(characterId)
		if s == nil {
			l.WithError(err).Errorf("Cannot locate session for instruction")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = s.Announce(writer.WriteHint(l)(cs.Data.Attributes.Message, cs.Data.Attributes.Width, cs.Data.Attributes.Height))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}
		err = s.Announce(writer.WriteEnableActions(l))
		if err != nil {
			l.WithError(err).Errorf("Unable to announce to character %d", s.CharacterId())
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func getCharacterId(l logrus.FieldLogger) func(r *http.Request) uint32 {
	return func(r *http.Request) uint32 {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["characterId"])
		if err != nil {
			l.Println("Error parsing characterId as uint32")
			return 0
		}
		return uint32(value)
	}
}
