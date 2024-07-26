package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-fonts/database"
	"net/http"
)

func getFontDesigners(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	font, err := database.GetDesigners(fontName)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	data, err := json.Marshal(font)
	if err != nil {
		http.Error(w, "Failed to load font designers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}
