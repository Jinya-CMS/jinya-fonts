package config

import (
	"fmt"
	"go-simpler.org/env"
)

type Configuration struct {
	ApiKey            string `env:"GOOGLE_API_KEY"`
	MongoUrl          string `env:"MONGO_URL"`
	MongoDatabase     string `env:"MONGO_DATABASE"`
	RedisUrl          string `env:"REDIS_URL"`
	ServeWebsite      bool   `env:"SERVE_WEBSITE"`
	OpenIDClientId    string `env:"OPENID_CLIENT_ID"`
	OpenIDDomain      string `env:"OPENID_DOMAIN"`
	OpenIDKeyFileData string `env:"OPENID_KEY_FILE_DATA"`
	ServerUrl         string `env:"SERVER_URL"`
}

func (c Configuration) GetRedirectUrl() string {
	return fmt.Sprintf("%s/admin/login/callback", c.ServerUrl)
}

var LoadedConfiguration *Configuration

func LoadConfiguration() error {
	config := new(Configuration)
	err := env.Load(config, nil)
	if err != nil {
		return err
	}

	LoadedConfiguration = config

	return nil
}
