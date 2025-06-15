package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-fonts/database"
	"net/http"
)

func GetFontDesigners(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	designers, err := database.GetDesigners(fontName)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(designers)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func createFontDesigner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]

	var designer database.Designer
	err := json.NewDecoder(r.Body).Decode(&designer)
	if err != nil {
		http.Error(w, "Failed to parse body", http.StatusBadRequest)
		return
	}

	newDesigner, err := database.CreateDesigner(name, designer)
	if err != nil {
		http.Error(w, "Failed to save designer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(newDesigner)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
		return
	}
}

func deleteFontDesigner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]
	designerName := vars["designerName"]

	err := database.DeleteDesigner(name, designerName)
	if err != nil {
		http.Error(w, "Failed to save designer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
