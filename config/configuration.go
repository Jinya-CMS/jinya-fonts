package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Configuration struct {
	ApiKey         string   `yaml:"api_key"`
	FontFileFolder string   `yaml:"font_file_folder"`
	FilterByName   []string `yaml:"filter_by_name"`
}

var LoadedConfiguration *Configuration

func LoadConfiguration(path string) (*Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	config := new(Configuration)

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)

	LoadedConfiguration = config

	return config, err
}
