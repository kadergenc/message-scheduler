package config

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"
)

var defaultRemoteServiceTimeout = 30000 // in ms

type PostgresConfig struct {
	WriteHost string `json:"writeHost"`
	WritePort string `json:"writePort"`
	ReadHost  string `json:"readHost"`
	ReadPort  string `json:"readPort"`
	User      string `json:"username"`
	Password  string `json:"password"`
	DbName    string `json:"dbname"`
}

type WebhookConfiguration struct {
	Host    string `json:"host"`
	Timeout int    `json:"timeout"`
}

type AppConfig struct {
	WebhookConfig WebhookConfiguration `json:"webhook"`
	Port          string               `json:"port"`
	AppName       string               `json:"appName"`
	TeamName      string               `json:"teamName"`
	Postgres      PostgresConfig       `json:"postgres"`
}

func Read() AppConfig {
	pgSecretsFileName := flag.String("pg", "./config/pg-credentials.json", "Path to PostgreSQL credentials file")
	configsFileName := flag.String("config", "./config/config.json", "Path to main config file")

	flag.Parse()

	var appCfg AppConfig

	pgSecrets, err := os.Open(*pgSecretsFileName)
	if err != nil {
		log.Fatalf("error occurred while opening pgSecrets %s", err)
	}

	var pgConf PostgresConfig
	secretsBytes, _ := io.ReadAll(pgSecrets)
	if err := json.Unmarshal(secretsBytes, &pgConf); err != nil {
		log.Fatalf("error occurred while unmarshalling pgSecrets %s", err)
	}

	configs, err := os.Open(*configsFileName)
	if err != nil {
		log.Fatalf("error occurred while opening configs %s", err)
	}

	configsBytes, _ := io.ReadAll(configs)
	if err := json.Unmarshal(configsBytes, &appCfg); err != nil {
		log.Fatalf("error occurred while unmarshalling configs %s", err)
	}

	appCfg.Postgres.User = pgConf.User
	appCfg.Postgres.Password = pgConf.Password

	appCfg.SetDefaults()

	return appCfg
}

func (appCfg *AppConfig) SetDefaults() {
	if appCfg.WebhookConfig.Timeout == 0 {
		appCfg.WebhookConfig.Timeout = defaultRemoteServiceTimeout
	}

}
