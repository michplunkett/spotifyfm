package config

import (
	"fmt"
	"os"
	"time"

	"github.com/michplunkett/spotifyfm/models"
)

// EnvironmentVariableController will handle creating and handling all fo the configs
type EnvVars interface {
	GetLastFMApiKey() string
	GetLastFMSharedSecret() string
	GetLastFMSessionKey() string
	SetLastFMSessionKey(sessionKey string)
	GetSpotifyClientID() string
	GetSpotifyClientSecret() string
	GetSpotifyUserName() string
	SetSpotifyInfo(token, refreshToken string, expirationTime time.Time)
	GetSpotifyToken() string
	SetSpotifyToken(token string)
	GetSpotifyTokenExpiration() time.Time
	SetSpotifyTokenExpiration(expirationTime time.Time)
	GetSpotifyRefreshToken() string
	SetSpotifyRefreshToken(refreshToken string)
}

type envVars struct {
	lastfmApiKey        string
	lastfmSharedSecret  string
	spotifyClientID     string
	spotifyClientSecret string
	spotifyUserName     string
	fileConfigs         *models.FileConfigs
}

func NewEnvVars() EnvVars {
	e := envVars{}

	e.fileConfigs = models.GetConfigValues()

	lastFmEnvVars := []string{os.Getenv(LastFMApiKey), os.Getenv(LastFMSharedSecret)}
	if arrayHasNoEmptyStrings(lastFmEnvVars) {
		e.lastfmApiKey = lastFmEnvVars[0]
		e.lastfmSharedSecret = lastFmEnvVars[1]
	} else {
		fmt.Errorf("one of the last.fm environment variables is not present in your system")
	}

	spotifyEnvVars := []string{os.Getenv(SpotifyClientID), os.Getenv(SpotifyClientSecret), os.Getenv(SpotifyUserName)}
	if arrayHasNoEmptyStrings(spotifyEnvVars) {
		e.spotifyClientID = spotifyEnvVars[0]
		e.spotifyClientSecret = spotifyEnvVars[1]
		e.spotifyUserName = spotifyEnvVars[2]
	} else {
		fmt.Errorf("one of the spotify environment variables is not present in your system")
	}

	return &e
}

func (e *envVars) GetLastFMApiKey() string {
	return e.lastfmApiKey
}

func (e *envVars) GetLastFMSharedSecret() string {
	return e.lastfmSharedSecret
}

func (e *envVars) GetLastFMSessionKey() string {
	return e.fileConfigs.LastFM.SessionKey
}

func (e *envVars) SetLastFMSessionKey(sessionKey string) {
	e.fileConfigs.LastFM.SessionKey = sessionKey
	e.fileConfigs.SetConfigValues()
}

func (e *envVars) GetSpotifyClientID() string {
	return e.spotifyClientID
}

func (e *envVars) GetSpotifyClientSecret() string {
	return e.spotifyClientSecret
}

func (e *envVars) GetSpotifyUserName() string {
	return e.spotifyUserName
}

func (e *envVars) GetSpotifyToken() string {
	return e.fileConfigs.Spotify.Token
}

func (e *envVars) SetSpotifyInfo(token, refreshToken string, expirationTime time.Time) {
	e.fileConfigs.Spotify.Token = token
	e.fileConfigs.Spotify.RefreshToken = refreshToken
	e.fileConfigs.Spotify.TokenExpirationTime = expirationTime
	e.fileConfigs.SetConfigValues()
}

func (e *envVars) SetSpotifyToken(token string) {
	e.fileConfigs.Spotify.Token = token
	e.fileConfigs.SetConfigValues()
}

func (e *envVars) GetSpotifyTokenExpiration() time.Time {
	return e.fileConfigs.Spotify.TokenExpirationTime
}

func (e *envVars) SetSpotifyTokenExpiration(expirationTime time.Time) {
	e.fileConfigs.Spotify.TokenExpirationTime = expirationTime
	e.fileConfigs.SetConfigValues()
}

func (e *envVars) GetSpotifyRefreshToken() string {
	return e.fileConfigs.Spotify.RefreshToken
}

func (e *envVars) SetSpotifyRefreshToken(refreshToken string) {
	e.fileConfigs.Spotify.RefreshToken = refreshToken
	e.fileConfigs.SetConfigValues()
}

func arrayHasNoEmptyStrings(envVars []string) bool {
	for _, value := range envVars {
		if value == EmptyString {
			return false
		}
	}

	return true
}
