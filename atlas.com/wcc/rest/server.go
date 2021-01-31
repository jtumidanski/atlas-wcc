package rest

import (
   "atlas-wcc/rest/resources"
   "github.com/gorilla/mux"
   "log"
   "net/http"
   "os"
   "time"
)

type Server struct {
   l  *log.Logger
   hs *http.Server
}

func NewServer(l *log.Logger) *Server {
   router := mux.NewRouter().StrictSlash(true).PathPrefix("/ms/csrv").Subrouter()
   router.Use(commonHeader)

   i := resources.NewInstructionResource(l)
   iRouter := router.PathPrefix("/worlds/{worldId}/channels/{channelId}/characters/{characterId}/instructions").Subrouter()
   iRouter.HandleFunc("", i.CreateInstruction).Methods(http.MethodPost)

   hs := http.Server{
      Addr:         ":8080",
      Handler:      router,
      ErrorLog:     l,                 // set the logger for the server
      ReadTimeout:  5 * time.Second,   // max time to read request from the client
      WriteTimeout: 10 * time.Second,  // max time to write response to the client
      IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
   }
   return &Server{l, &hs}
}

func (s *Server) Run() {
   s.l.Println("[INFO] Starting server on port 8080")
   err := s.hs.ListenAndServe()
   if err != nil {
      s.l.Printf("Error starting server: %s\n", err)
      os.Exit(1)
   }
}

func commonHeader(next http.Handler) http.Handler {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      w.Header().Add("Content-Type", "application/json")
      next.ServeHTTP(w, r)
   })
}
