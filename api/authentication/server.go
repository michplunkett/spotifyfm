package authentication

import (
	"log"
	"net/http"
)

type Server interface {
	Start()
	Stop()
}

type authenticationServer struct {
	server *http.Server
}

func NewServer() Server {
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
	return &authenticationServer{
		server: server,
	}
}

func (authServer *authenticationServer) Start() {
	log.Fatal(authServer.server.ListenAndServe())
	log.Println("it has started")
}

func (authServer *authenticationServer) Stop() {
	authServer.server.Close()
}
