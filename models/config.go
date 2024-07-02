package models

import (
	"encoding/json"
	"os"
	"time"
)

const (
	configFileName = "config.json"
)

type FileConfigs struct {
	LastFM  LastFMConfigs
	Spotify SpotifyConfigs
}

type LastFMConfigs struct {
	SessionKey string
}

type SpotifyConfigs struct {
	Token               string
	RefreshToken        string
	TokenExpirationTime time.Time
}

func GetConfigValues() *FileConfigs {
	configFile, _ := os.ReadFile(configFileName)

	configInfo := FileConfigs{}
	_ = json.Unmarshal(configFile, &configInfo)

	return &configInfo
}

func (config *FileConfigs) SetConfigValues() {
	setConfigFileValues(config)
}

func setConfigFileValues(configValues *FileConfigs) {
	file, _ := json.MarshalIndent(configValues, "", " ")
	_ = os.WriteFile(configFileName, file, 0644)
}
