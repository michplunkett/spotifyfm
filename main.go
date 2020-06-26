package main

import (
	"fmt"

	"github.com/michplunkett/spotifyfm/config"
)

func main() {
	envVarCtrl := config.NewEnvVarController()
	envVarCtrl.Init()
	fmt.Println(envVarCtrl.GetLastFMConfig())
	fmt.Println(envVarCtrl.GetSpotifyConfig())
}
