package models

import (
	"encoding/json"
	"io/ioutil"
)

var (
	ConfigFileName = "config.json"
)

type FileConfigs struct {
	LastFM  LastFMConfigs
	Spotify SpotifyConfigs
}

type LastFMConfigs struct {
	SessionKey string
}

type SpotifyConfigs struct {
	StartTime           string
	Token               string
	RefreshToken        string
	TokenExpirationTime string
}

func GetConfigValues() *FileConfigs {
	configFile, _ := ioutil.ReadFile(ConfigFileName)

	configInfo := FileConfigs{}
	json.Unmarshal(configFile, &configInfo)

	return &configInfo
}

func (config *FileConfigs) SetConfigValues() {
	setConfigFileValues(config)
}

func setConfigFileValues(configValues *FileConfigs) {
	file, _ := json.MarshalIndent(configValues, "", " ")
	_ = ioutil.WriteFile(ConfigFileName, file, 0644)
}