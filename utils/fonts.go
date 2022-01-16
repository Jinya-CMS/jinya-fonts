package utils

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
)

func LoadFonts() ([]meta.FontFile, error) {
	files, err := ioutil.ReadDir(config.LoadedConfiguration.FontFileFolder)
	if err != nil {
		return nil, err
	}

	var availableFonts []meta.FontFile

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		yamlFileData, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + file.Name())
		if err != nil {
			continue
		}

		var fontFile meta.FontFile
		err = yaml.Unmarshal(yamlFileData, &fontFile)
		if err != nil {
			continue
		}

		availableFonts = append(availableFonts, fontFile)
	}

	return availableFonts, err
}
