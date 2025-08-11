package api

import (
	"encoding/json"
	"jinya-fonts/database"
	"jinya-fonts/fontsync"
	"log"
	"net/http"

	"github.com/bamzi/jobrunner"
)

func getSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := database.GetSettings()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(settings)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
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

	for _, entry := range jobrunner.Entries() {
		jobrunner.Remove(entry.ID)
	}

	if settings.SyncEnabled {
		err = jobrunner.Schedule(settings.SyncInterval, fontsync.SyncJob{})
		if err != nil {
			log.Printf("Failed to schedule sync job %s", err.Error())
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
