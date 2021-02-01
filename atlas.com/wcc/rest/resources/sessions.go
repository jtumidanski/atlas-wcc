package resources

import (
	"atlas-wcc/mapleSession"
	"atlas-wcc/registries"
	"atlas-wcc/rest/attributes"
	"log"
	"net/http"
	"strconv"
)

type SessionListDataContainer struct {
	Data []SessionData `json:"data"`
}

type SessionData struct {
	Id         string            `json:"id"`
	Type       string            `json:"type"`
	Attributes SessionAttributes `json:"attributes"`
}

type SessionAttributes struct {
	AccountId uint32 `json:"accountId"`
	WorldId   byte   `json:"worldId"`
	ChannelId byte   `json:"channelId"`
	CharacterId uint32 `json:"characterId"`
}

type SessionResource struct {
	l *log.Logger
}

func NewSessionResource(l *log.Logger) *SessionResource {
	return &SessionResource{l}
}

func (s *SessionResource) GetSessions(rw http.ResponseWriter, _ *http.Request) {
	ss := registries.GetSessionRegistry().GetAll()

	var response SessionListDataContainer
	response.Data = make([]SessionData, 0)
	for _, x := range ss {
		response.Data = append(response.Data, *getSessionObject(x))
	}

	err := attributes.ToJSON(response, rw)
	if err != nil {
		s.l.Println("Error encoding GetSessions response")
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func getSessionObject(x mapleSession.MapleSession) *SessionData {
	return &SessionData{
		Id:   strconv.Itoa(x.SessionId()),
		Type: "Session",
		Attributes: SessionAttributes{
			AccountId: x.AccountId(),
			WorldId:   x.WorldId(),
			ChannelId: x.ChannelId(),
			CharacterId: x.CharacterId(),
		},
	}
}
