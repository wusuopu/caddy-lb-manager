package config

import (
	"os"
)

type IConfig struct {
	Server struct {
		Port		string
		GO_ENV	string
		BaseUrl	string
		Version	string
	}
}

var Config IConfig

func Load() IConfig {
	Config.Server.Port = os.Getenv("PORT")
	if Config.Server.Port == "" {
		Config.Server.Port = "80"
	}

	Config.Server.GO_ENV = os.Getenv("GO_ENV")
	if Config.Server.GO_ENV == "" {
		Config.Server.GO_ENV = "development"
	}

	Config.Server.Version = "0.1.0"

	Config.Server.BaseUrl = os.Getenv("WEBUI_BASE_URL")
	if Config.Server.BaseUrl == "/" {
		Config.Server.BaseUrl = ""
	}
	return Config
}