package api

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"jinya-fonts/database"
	"net/http"
)

func getFontFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]

	files, err := database.GetFontFiles(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func createFontFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]
	weight := vars["fontWeight"]
	style := vars["fontStyle"]
	fontType := vars["fontType"]
	if fontType != "woff2" && fontType != "ttf" {
		http.Error(w, "Invalid font type", http.StatusBadRequest)
		return
	}

	fileBuffer := bytes.NewBufferString("")
	_, err := io.Copy(fileBuffer, r.Body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}

	_, err = database.SetFontFile(name, weight, style, fontType, fileBuffer.Bytes())
	if err != nil {
		http.Error(w, "Failed to create font file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateFontFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]
	weight := vars["fontWeight"]
	style := vars["fontStyle"]
	fontType := vars["fontType"]
	if fontType != "woff2" && fontType != "ttf" {
		http.Error(w, "Invalid font type", http.StatusBadRequest)
		return
	}

	fileBuffer := bytes.NewBufferString("")
	_, err := io.Copy(fileBuffer, r.Body)
	if err != nil {
		http.Error(w, "Could not read body", http.StatusBadRequest)
		return
	}

	_, err = database.SetFontFile(name, weight, style, fontType, fileBuffer.Bytes())
	if err != nil {
		return
	}
	if err != nil {
		http.Error(w, "Failed to update font file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteFontFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["fontName"]
	weight := vars["fontWeight"]
	style := vars["fontStyle"]
	fontType := vars["fontType"]
	if fontType != "woff2" && fontType != "ttf" {
		http.Error(w, "Invalid font type", http.StatusBadRequest)
		return
	}

	err := database.RemoveFontFile(name, weight, style, fontType)
	if err != nil {
		http.Error(w, "Failed to delete font file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
