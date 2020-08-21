package config

import (
	"fmt"
	"os"

	"github.com/michplunkett/spotifyfm/utility"
)

// EnvironmentVariableController will handle creating and handling all fo the configs
type EnvVarController interface {
	Init()
	GetLastFMConfig() *LastFmConfig
	GetSpotifyConfig() *SpotifyConfig
}

type LastFmConfig struct {
	apiKey       string
	sharedSecret string
}

type SpotifyConfig struct {
	clientID     string
	clientSecret string
	userName     string
}

type envVarController struct {
	lastFmConfigObj  *LastFmConfig
	spotifyConfigObj *SpotifyConfig
}

func NewEnvVarController() EnvVarController {
	return &envVarController{}
}

// Init creates the two variables that will be exported in the GetLastFMConfig and GetSpotifyConfig functions
func (e *envVarController) Init() {
	lastFmEnvVars := []string{os.Getenv(utility.LastFMApiKey), os.Getenv(utility.LastFMSharedSecret)}
	if utility.ArrayHasNoEmptyStrings(lastFmEnvVars) {
		e.lastFmConfigObj = &LastFmConfig{
			apiKey:       lastFmEnvVars[0],
			sharedSecret: lastFmEnvVars[1],
		}
	} else {
		fmt.Errorf("one of the last.fm environment variables is not present in your system")
	}

	spotifyEnvVars := []string{os.Getenv(utility.SpotifyClientID), os.Getenv(utility.SpotifyClientSecret), os.Getenv(utility.SpotifyUserName)}
	if utility.ArrayHasNoEmptyStrings(spotifyEnvVars) {
		e.spotifyConfigObj = &SpotifyConfig{
			clientID:     spotifyEnvVars[0],
			clientSecret: spotifyEnvVars[1],
			userName:     spotifyEnvVars[2],
		}
	} else {
		fmt.Errorf("one of the spotify environment variables is not present in your system")
	}
}

func (e *envVarController) GetLastFMConfig() *LastFmConfig {
	return e.lastFmConfigObj
}

func (e *envVarController) GetSpotifyConfig() *SpotifyConfig {
	return e.spotifyConfigObj
}

func (config *LastFmConfig) GetApiKey() string {
	return config.apiKey
}

func (config *LastFmConfig) GetSharedSecret() string {
	return config.sharedSecret
}

func (config *SpotifyConfig) GetClientID() string {
	return config.clientID
}

func (config *SpotifyConfig) GetClientSecret() string {
	return config.clientSecret
}

func (config *SpotifyConfig) GetUserName() string {
	return config.userName
}
