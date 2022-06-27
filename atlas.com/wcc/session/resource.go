package session

import (
	"atlas-wcc/json"
	"atlas-wcc/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

const (
	GetSessions = "get_sessions"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	sRouter := router.PathPrefix("/sessions").Subrouter()
	sRouter.HandleFunc("", registerGetSessions(l)).Methods(http.MethodGet)
}

func registerGetSessions(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(GetSessions, func(span opentracing.Span) http.HandlerFunc {
		return handleGetSessions(l)(span)
	})
}

func handleGetSessions(l logrus.FieldLogger) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			ss := Registry().GetAll()
			var response DataListContainer
			response.Data = make([]DataBody, 0)
			for _, x := range ss {
				response.Data = append(response.Data, getSessionObject(x))
			}

			err := json.ToJSON(response, w)
			if err != nil {
				l.WithError(err).Errorf("Encoding response")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
}

func getSessionObject(x Model) DataBody {
	return DataBody{
		Id:   strconv.Itoa(int(x.SessionId())),
		Type: "Session",
		Attributes: Attributes{
			AccountId:   x.AccountId(),
			WorldId:     x.WorldId(),
			ChannelId:   x.ChannelId(),
			CharacterId: x.CharacterId(),
		},
	}
}
