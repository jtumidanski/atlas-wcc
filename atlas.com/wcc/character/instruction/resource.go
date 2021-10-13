package instruction

import (
	"atlas-wcc/json"
	"atlas-wcc/rest"
	"atlas-wcc/session"
	"atlas-wcc/socket/response/writer"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	CreateInstruction = "create_instruction"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	iRouter := router.PathPrefix("/characters/{characterId}/instructions").Subrouter()
	iRouter.HandleFunc("", registerCreateInstruction(l)).Methods(http.MethodPost)
}

func registerCreateInstruction(l logrus.FieldLogger) func(http.ResponseWriter, *http.Request) {
	return rest.RetrieveSpan(CreateInstruction, func(span opentracing.Span) http.HandlerFunc {
		return ParseCharacterId(l, func(characterId uint32) http.HandlerFunc {
			return ParseInput(l, func(input *InputDataContainer) http.HandlerFunc {
				return handleCreateInstruction(l)(span)(characterId)(input)
			})
		})
	})
}

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

type CharacterIdHandler func(characterId uint32) http.HandlerFunc

func ParseCharacterId(l logrus.FieldLogger, next CharacterIdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["characterId"])
		if err != nil {
			l.WithError(err).Errorln("Error parsing id as uint32")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(value))(w, r)
	}
}

type InputHandler func(input *InputDataContainer) http.HandlerFunc

func ParseInput(l logrus.FieldLogger, next InputHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		i := &InputDataContainer{}
		err := json.FromJSON(i, r.Body)
		if err != nil {
			l.WithError(err).Errorf("Deserializing instruction")
			w.WriteHeader(http.StatusBadRequest)
			err := json.ToJSON(&GenericError{Message: err.Error()}, w)
			if err != nil {
				l.WithError(err).Errorf("Unable to serialize error mesage")
			}
			return
		}
		next(i)(w, r)
	}
}

func handleCreateInstruction(l logrus.FieldLogger) func(span opentracing.Span) func(characterId uint32) func(input *InputDataContainer) http.HandlerFunc {
	return func(span opentracing.Span) func(characterId uint32) func(input *InputDataContainer) http.HandlerFunc {
		return func(characterId uint32) func(input *InputDataContainer) http.HandlerFunc {
			return func(input *InputDataContainer) http.HandlerFunc {
				return func(w http.ResponseWriter, r *http.Request) {
					s := session.GetByCharacterId(characterId)
					if s == nil {
						l.Errorf("Cannot locate session for instruction")
						w.WriteHeader(http.StatusBadRequest)
						return
					}

					attr := input.Data.Attributes
					err := s.Announce(writer.WriteHint(l)(attr.Message, attr.Width, attr.Height))
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
		}
	}
}
