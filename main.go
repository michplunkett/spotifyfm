package main

import (
	"log"
	"net/http"

	"github.com/michplunkett/spotifyfm/api/authentication"
	"github.com/michplunkett/spotifyfm/config"
)

func main() {
	// Start up http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)

	spotifyProfileHandler := authentication.NewSpotifyAuthHandlerProfile()
	spotifyProfileHandler.Authenticate()

	lastFMConfig := config.NewEnvVarController()
	lastFMConfig.Init()
	lastFMHandler := authentication.NewLastFMAuthHandler(lastFMConfig.GetLastFMConfig())
	lastFMHandler.Authenticate()


}
