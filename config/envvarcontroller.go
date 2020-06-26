package config

import (
	"fmt"
	"os"

	"github.com/michplunkett/spotifyfm/utility"
	)

// EnvironmentVariableController will handle creating and handling all fo the configs
type EnvVarController interface {
	Init() error
	GetLastFMConfig() (*LastFmConfig, error)
	GetSpotifyConfig() (*SpotifyConfig, error)
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
	lastFmConfigObj *LastFmConfig
	spotifyConfigObj *SpotifyConfig
}

func NewEnvVarController() EnvVarController {
	return &envVarController{}
}

// Init creates the two variables that will be exported in the GetLastFMConfig and GetSpotifyConfig functions
func (e *envVarController) Init() error {
	lastFmEnvVars := []string{os.Getenv(utility.LastFMApiKey), os.Getenv(utility.LastFMSharedSecret)}
	if utility.ArrayHasNoEmptyStrings(lastFmEnvVars) {
		return fmt.Errorf("one of the last.fm environment variables is not present in your system")
	}

	e.lastFmConfigObj = &LastFmConfig{
		apiKey:       lastFmEnvVars[0],
		sharedSecret: lastFmEnvVars[1],
	}

	spotifyEnvVars := []string{os.Getenv(utility.SpotifyClientID), os.Getenv(utility.SpotifyClientSecret), os.Getenv(utility.SpotifyUserName)}
	if utility.ArrayHasNoEmptyStrings(spotifyEnvVars) {
		return fmt.Errorf("one of the spotify environment variables is not present in your system")
	}


	e.spotifyConfigObj = &SpotifyConfig{
		clientID:     spotifyEnvVars[0],
		clientSecret: spotifyEnvVars[1],
		userName:     spotifyEnvVars[2],
	}

	return nil
}

func (e *envVarController) GetLastFMConfig() (*LastFmConfig, error) {
	return e.lastFmConfigObj, nil
}

func (e *envVarController) GetSpotifyConfig() (*SpotifyConfig, error) {
	return e.spotifyConfigObj, nil
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
