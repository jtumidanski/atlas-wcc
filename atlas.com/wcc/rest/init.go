package rest

import (
	"atlas-wcc/character/instruction"
	"atlas-wcc/session"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync"
)

func CreateRestService(l *logrus.Logger, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg, ProduceRoutes)
}

func ProduceRoutes(l logrus.FieldLogger) http.Handler {
	router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/csrv/worlds/{worldId}/channels/{channelId}").Subrouter()
	router.Use(CommonHeader)

	sRouter := router.PathPrefix("/sessions").Subrouter()
	sRouter.HandleFunc("", session.HandleGetSessions(l)).Methods(http.MethodGet)

	iRouter := router.PathPrefix("/characters/{characterId}/instructions").Subrouter()
	iRouter.HandleFunc("", instruction.HandleCreateInstruction(l)).Methods(http.MethodPost)

	return router
}
