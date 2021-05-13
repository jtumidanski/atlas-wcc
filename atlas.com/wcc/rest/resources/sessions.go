package resources

import (
	"atlas-wcc/json"
	"atlas-wcc/mapleSession"
	"atlas-wcc/registries"
	"atlas-wcc/rest/attributes"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type SessionResource struct {
	l logrus.FieldLogger
}

func NewSessionResource(l logrus.FieldLogger) *SessionResource {
	return &SessionResource{l}
}

func (s *SessionResource) GetSessions(rw http.ResponseWriter, _ *http.Request) {
	ss := registries.GetSessionRegistry().GetAll()

	var response attributes.SessionListDataContainer
	response.Data = make([]attributes.SessionData, 0)
	for _, x := range ss {
		response.Data = append(response.Data, *getSessionObject(x))
	}

	err := json.ToJSON(response, rw)
	if err != nil {
		s.l.WithError(err).Errorf("Encoding GetSessions response")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func getSessionObject(x mapleSession.MapleSession) *attributes.SessionData {
	return &attributes.SessionData{
		Id:   strconv.Itoa(x.SessionId()),
		Type: "Session",
		Attributes: attributes.SessionAttributes{
			AccountId:   x.AccountId(),
			WorldId:     x.WorldId(),
			ChannelId:   x.ChannelId(),
			CharacterId: x.CharacterId(),
		},
	}
}
