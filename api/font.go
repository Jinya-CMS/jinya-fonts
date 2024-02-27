package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-fonts/database"
	"net/http"
)

type addFontData struct {
	Name        string `json:"name"`
	License     string `json:"license"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type updateFontData struct {
	Name        string `json:"name"`
	License     string `json:"license"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func getAllFonts(w http.ResponseWriter, r *http.Request) {
	availableFonts, err := database.GetAllFonts()
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
	_, _ = w.Write(data)
}

func getFontByName(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	font, err := database.GetFont(fontName)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	data, err := json.Marshal(font)
	if err != nil {
		http.Error(w, "Failed to load font", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
}

func createFont(w http.ResponseWriter, r *http.Request) {
	font := new(addFontData)

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(font)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	webfont, err := database.CreateFont(font.Name, font.License, font.Category, font.Description)
	if err != nil {
		http.Error(w, "Failed to create font", http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	w.WriteHeader(http.StatusCreated)
	_ = encoder.Encode(webfont)
}

func updateFont(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	_, err := database.GetFont(fontName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	font := new(updateFontData)

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(font)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = database.UpdateFont(fontName, font.Description, font.License, font.Category)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFont(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	_, err := database.GetFont(fontName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = database.DeleteFont(fontName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
