package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-fonts/database"
	"net/http"
)

func getFontFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]

	fonts, err := database.GetFontFiles(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	data, err := json.Marshal(fonts)
	if err != nil {
		http.Error(w, "Failed to load font", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
