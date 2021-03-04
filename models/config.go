package models

import (
	"encoding/json"
	"io/ioutil"
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
	configFile, _ := ioutil.ReadFile(configFileName)

	configInfo := FileConfigs{}
	_ = json.Unmarshal(configFile, &configInfo)

	return &configInfo
}

func (config *FileConfigs) SetConfigValues() {
	setConfigFileValues(config)
}

func setConfigFileValues(configValues *FileConfigs) {
	file, _ := json.MarshalIndent(configValues, "", " ")
	_ = ioutil.WriteFile(configFileName, file, 0644)
}
