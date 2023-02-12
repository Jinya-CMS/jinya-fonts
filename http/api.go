package http

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"jinya-fonts/utils"
	"net/http"
	"os"
)

func GetFontMeta(w http.ResponseWriter, r *http.Request) {
	font := r.URL.Query().Get("font")
	if font != "" {
		yamlFileData, err := os.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + font + ".yaml")
		var fontFile meta.FontFile
		err = yaml.Unmarshal(yamlFileData, &fontFile)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		data, err := json.Marshal(fontFile)
		if err != nil {
			http.Error(w, "Failed to load font", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		availableFonts, err := utils.LoadFonts()
		if err != nil {
			http.NotFound(w, r)
			return
		}

		data, err := json.Marshal(availableFonts)
		if err != nil {
			http.Error(w, "Failed to load fonts", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	}
}
