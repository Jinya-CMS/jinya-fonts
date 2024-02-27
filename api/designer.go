package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-fonts/database"
	"net/http"
)

type addDesigner struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

type editDesigner struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

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

func createFontDesigner(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]

	designer := addDesigner{}
	err := json.NewDecoder(r.Body).Decode(&designer)
	if err != nil {
		http.Error(w, "Failed to parse body", http.StatusBadRequest)
		return
	}

	newDesigner, err := database.CreateDesigner(name, designer.Name, designer.Bio)
	if err != nil {
		http.Error(w, "Failed to save designer", http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(newDesigner)
	if err != nil {
		http.Error(w, "Failed to save designer", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
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
