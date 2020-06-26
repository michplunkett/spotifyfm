package apis

import (
	"github.com/michplunkett/spotifyfm/config"
	"github.com/shkh/lastfm-go/lastfm"
)

type LastFMAPI interface {
	// Gonna fill this up with lastfm info grabbing
}

type lastFMAPI struct {
	api *lastfm.Api
}

func NewLastFMAPI(config *config.LastFmConfig) LastFMAPI {
	api := lastfm.New(config.GetApiKey(), config.GetSharedSecret())
	return &lastFMAPI{api}
}
