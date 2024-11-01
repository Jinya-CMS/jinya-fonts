package config

import (
	"fmt"
	"go-simpler.org/env"
)

type Configuration struct {
	ApiKey                 string `env:"GOOGLE_API_KEY"`
	MongoUrl               string `env:"MONGO_URL"`
	MongoDatabase          string `env:"MONGO_DATABASE"`
	RedisUrl               string `env:"REDIS_URL"`
	ServeWebsite           bool   `env:"SERVE_WEBSITE"`
	OidcFrontendClientId   string `env:"OIDC_FRONTEND_CLIENT_ID"`
	OidcDomain             string `env:"OIDC_DOMAIN"`
	OidcServerClientId     string `env:"OIDC_SERVER_CLIENT_ID"`
	OidcServerClientSecret string `env:"OIDC_SERVER_CLIENT_SECRET"`
	ServerUrl              string `env:"SERVER_URL"`
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
