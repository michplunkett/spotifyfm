package constants

import "time"

const (
	// General
	APIObjectLimit     = 200
	AverageSongSeconds = 210
	DoubleHyphen       = "--"
	EmptyString        = ""
	NotFound           = "NF"

	// Last.fm credential names
	LastFMApiKey       = "LAST_FM_API_KEY"
	LastFMSharedSecret = "LAST_FM_SHARED_SECRET"
	LastFMRedirectURL  = "http://localhost:8080/lastfm-callback"

	// Spotify credential names
	SpotifyClientID     = "SPOTIFY_ID"
	SpotifyClientSecret = "SPOTIFY_SECRET"
	SpotifyRedirectURL  = "http://localhost:8080/spotify-callback"
	SpotifyUserName     = "SPOTIFY_USER_NAME"

	// Time durations
	LastFMPeriod7Day    = "7day"
	LastFMPeriod1Month  = "1month"
	LastFMPeriod3Month  = "3month"
	LastFMPeriod6Month  = "6month"
	LastFMPeriod12Month = "12month"
	SpotifyPeriodLong   = "long"
	SpotifyPeriodMedium = "medium"
	SpotifyPeriodShort  = "short"
)

var (
	// There are some artists that have aliases in last.fm
	ProblematicArtistMapping = map[string]string{
		"strfkr":   "starfucker",
		"th1rt3en": "thirteen",
	}
	Now              = time.Now()
	StartOf2021      = time.Date(2021, time.January, 1, 0, 0, 0, 0, Now.Location()).Unix()
	StartOf2020      = time.Date(2020, time.January, 1, 0, 0, 0, 0, Now.Location()).Unix()
	StartOfThisMonth = time.Date(Now.Year(), Now.Month(), 1, 0, 0, 0, 0, Now.Location()).Unix()
)
