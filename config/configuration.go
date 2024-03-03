package config

import (
	"go-simpler.org/env"
)

type Configuration struct {
	ApiKey         string `env:"GOOGLE_API_KEY"`
	MongoUrl       string `env:"MONGO_URL"`
	MongoDatabase  string `env:"MONGO_DATABASE"`
	GoogleRedisUrl string `env:"GOOGLE_REDIS_URL"`
	CustomRedisUrl string `env:"CUSTOM_REDIS_URL"`
	ServeWebsite   bool   `env:"SERVE_WEBSITE"`
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
