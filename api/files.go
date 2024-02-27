package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"jinya-fonts/database"
	"net/http"
)

type addFontFile struct {
	FontFileData string `json:"data"`
	Subset       string `json:"subset"`
	Weight       string `json:"weight"`
	Style        string `json:"style"`
}

type editFontFile struct {
	FontFileData string `json:"data"`
	Subset       string `json:"subset"`
	Weight       string `json:"weight"`
	Style        string `json:"style"`
}

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

func createFontFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]
	subset := vars["fontSubset"]
	weight := vars["fontWeight"]
	style := vars["fontStyle"]

	fileBuffer := bytes.NewBufferString("")
	_, err := io.Copy(fileBuffer, r.Body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}

	_, err = database.AddFontFile(name, fileBuffer.Bytes(), subset, weight, style)
	if err != nil {
		http.Error(w, "Failed to create font file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateFontFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]
	subset := vars["fontSubset"]
	weight := vars["fontWeight"]
	style := vars["fontStyle"]

	fileBuffer := bytes.NewBufferString("")
	_, err := io.Copy(fileBuffer, r.Body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}

	err = database.UpdateFontFile(name, fileBuffer.Bytes(), subset, weight, style)
	if err != nil {
		http.Error(w, "Failed to update font file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFontFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]
	subset := vars["fontSubset"]
	weight := vars["fontWeight"]
	style := vars["fontStyle"]

	err := database.DeleteFontFile(name, subset, weight, style)
	if err != nil {
		http.Error(w, "Failed to delete font file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
