package resources

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/registries"
	"atlas-wcc/rest/attributes"
	"log"
	"net/http"
	"strconv"
)

type SessionResource struct {
	l *log.Logger
}

func NewSessionResource(l *log.Logger) *SessionResource {
	return &SessionResource{l}
}

func (s *SessionResource) GetSessions(rw http.ResponseWriter, _ *http.Request) {
	ss := registries.GetSessionRegistry().GetAll()

	var response attributes.SessionListDataContainer
	response.Data = make([]attributes.SessionData, 0)
	for _, x := range ss {
		response.Data = append(response.Data, *getSessionObject(x))
	}

	err := attributes.ToJSON(response, rw)
	if err != nil {
		s.l.Println("Error encoding GetSessions response")
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
