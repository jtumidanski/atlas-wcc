package rest

import (
	"atlas-wcc/rest/resources"
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

	s := resources.NewSessionResource(l)
	sRouter := router.PathPrefix("/sessions").Subrouter()
	sRouter.HandleFunc("", s.GetSessions)

	i := resources.NewInstructionResource(l)
	iRouter := router.PathPrefix("/characters/{characterId}/instructions").Subrouter()
	iRouter.HandleFunc("", i.CreateInstruction).Methods(http.MethodPost)

	return router
}
