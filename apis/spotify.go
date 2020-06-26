package apis

import "github.com/michplunkett/spotifyfm/config"

type SpotifyAPI interface {
	// Gonna flesh this out with spotify info grabbing
}

type spotifyAPI struct {
	spotifyConfig *config.SpotifyConfig
}

func NewSpotifyAPI(config *config.SpotifyConfig) SpotifyAPI {
	// Next step is to set up spotify authentication: https://github.com/zmb3/spotify#authentication
	return &spotifyAPI{config}
}
