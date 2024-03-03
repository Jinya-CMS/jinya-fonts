package api

import (
	"encoding/json"
	"jinya-fonts/database"
	"net/http"
)

func getSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := database.GetSettings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(settings)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func updateSettings(w http.ResponseWriter, r *http.Request) {
	settings := new(database.JinyaFontsSettings)
	err := json.NewDecoder(r.Body).Decode(settings)
	if err != nil {
		http.Error(w, "Failed to parse body", http.StatusBadRequest)
		return
	}

	err = database.UpdateSettings(settings)
	if err != nil {
		http.Error(w, "Failed to save settings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
