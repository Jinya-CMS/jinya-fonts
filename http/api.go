package http

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
)

func GetFontList(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(config.LoadedConfiguration.FontFileFolder)
	if err != nil {
		http.NotFound(w, r)
		return
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

	data, err := json.Marshal(availableFonts)
	if err != nil {
		http.Error(w, "Failed to load fonts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
