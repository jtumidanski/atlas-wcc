package session

import (
	"atlas-wcc/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func HandleGetSessions(l logrus.FieldLogger) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		ss := GetRegistry().GetAll()

		var response DataListContainer
		response.Data = make([]DataBody, 0)
		for _, x := range ss {
			response.Data = append(response.Data, getSessionObject(x))
		}

		err := json.ToJSON(response, w)
		if err != nil {
			l.WithError(err).Errorf("Encoding GetSessions response")
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getSessionObject(x *Model) DataBody {
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
