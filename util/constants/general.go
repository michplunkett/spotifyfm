package constants

const (
	// General
	EmptyString    = ""
	APIObjectLimit = 300

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

// There are some artists that have aliases in last.fm
var ProblematicArtistMapping = map[string]string{
	"strfkr": "starfucker",
	"th1rt3en": "thirteen",
}
