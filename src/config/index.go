package config

import (
	"os"
)

type IConfig struct {
	Server struct {
		Port		string
		GO_ENV	string
		ClusterMode bool
		BaseUrl	string
		Version	string
		Username string
		Password string
	}
	Caddy struct {
		BinPath		string
		DataPath	string
		ConfigPath	string
		ReloadCMD	string
		TLSEmail	string
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

	Config.Server.ClusterMode = os.Getenv("WEBUI_ENABLE_CLUSTER_MODE") == "true"

	Config.Server.Version = "0.1.0"

	Config.Server.BaseUrl = os.Getenv("WEBUI_BASE_URL")
	if Config.Server.BaseUrl == "/" {
		Config.Server.BaseUrl = ""
	}

	Config.Server.Username = os.Getenv("WEBUI_BASIC_AUTH_USER")
	Config.Server.Password = os.Getenv("WEBUI_BASIC_AUTH_PASSWORD")

	// Caddy 配置
	Config.Caddy.BinPath = os.Getenv("CADDY_BIN_PATH")
	if Config.Caddy.BinPath == "" {
		Config.Caddy.BinPath = "caddy"
	}

	Config.Caddy.DataPath = os.Getenv("CADDY_DATA_PATH")
	if Config.Caddy.DataPath == "" {
		Config.Caddy.DataPath = "/data/caddy"
	}

	Config.Caddy.ConfigPath = os.Getenv("CADDY_CONFIG_PATH")
	if Config.Caddy.ConfigPath == "" {
		Config.Caddy.ConfigPath = "/data/Caddyfile"
	}

	Config.Caddy.ReloadCMD = os.Getenv("CADDY_RELOAD_CMD")

	Config.Caddy.TLSEmail = os.Getenv("CADDY_TLS_EMAIL")

	return Config
}