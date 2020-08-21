package main

import "github.com/michplunkett/spotifyfm/api/authentication"

func main() {
	spotifyProfileHandler := authentication.NewSpotifyAuthHandlerProfile()
	spotifyProfileHandler.Authenticate()

	//lastFMConfig := config.NewEnvVarController()
	//lastFMConfig.Init()
	//lastFMHandler := authentication.NewLastFMAuthHandler(lastFMConfig.GetLastFMConfig())
	//lastFMHandler.Authenticate()
}
