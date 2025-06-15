package config

import (
	"fmt"
	"go-simpler.org/env"
	"os"
)

type Configuration struct {
	ApiKey                 string `env:"GOOGLE_API_KEY"`
	PostgresUrl            string `env:"POSTGRES_URL"`
	RedisUrl               string `env:"REDIS_URL"`
	ServeWebsite           bool   `env:"SERVE_WEBSITE"`
	OidcFrontendClientId   string `env:"OIDC_FRONTEND_CLIENT_ID"`
	OidcDomain             string `env:"OIDC_DOMAIN"`
	OidcServerClientId     string `env:"OIDC_SERVER_CLIENT_ID"`
	OidcServerClientSecret string `env:"OIDC_SERVER_CLIENT_SECRET"`
	ServerUrl              string `env:"SERVER_URL"`
	S3ServerUrl            string `env:"S3_SERVER_URL"`
	S3AccessKey            string `env:"S3_ACCESS_KEY"`
	S3SecretKey            string `env:"S3_SECRET_KEY"`
	S3Bucket               string `env:"S3_BUCKET"`
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

func IsDev() bool {
	return !IsProd()
}

func IsProd() bool {
	return os.Getenv("ENV") == "prod"
}
