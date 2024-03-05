package api

import (
	"cmp"
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-fonts/database"
	"jinya-fonts/fontsync"
	"net/http"
	"slices"
)

type addFontData struct {
	Name        string `json:"name"`
	License     string `json:"license"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type updateFontData struct {
	License     string `json:"license"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type apiFont struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Designers   []database.Designer `json:"designers"`
	License     string              `json:"license"`
	Category    string              `json:"category"`
	Fonts       []database.File     `json:"fonts"`
	GoogleFont  bool                `json:"googleFont"`
}

func convertWebfontToApiFont(webfont *database.Webfont) apiFont {
	fonts := make([]database.File, len(webfont.Fonts))
	i := 0
	for _, file := range webfont.Fonts {
		fonts[i] = file
		i += 1
	}

	slices.SortFunc(fonts, func(a, b database.File) int {
		return cmp.Compare(a.Path, b.Path)
	})

	font := apiFont{
		Fonts:       fonts,
		Name:        webfont.Name,
		Description: webfont.Description,
		Designers:   webfont.Designers,
		License:     webfont.License,
		Category:    webfont.Category,
		GoogleFont:  webfont.GoogleFont,
	}

	return font
}

func convertWebfontListToApiFontList(webfonts []database.Webfont) []apiFont {
	fonts := make([]apiFont, len(webfonts))
	for i, webfont := range webfonts {
		fonts[i] = convertWebfontToApiFont(&webfont)
	}

	return fonts
}

func getAllFonts(w http.ResponseWriter, r *http.Request) {
	availableFonts, err := database.GetAllFonts()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(convertWebfontListToApiFontList(availableFonts))
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func getGoogleFonts(w http.ResponseWriter, r *http.Request) {
	availableFonts, err := database.GetGoogleFonts()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(convertWebfontListToApiFontList(availableFonts))
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func getCustomFonts(w http.ResponseWriter, r *http.Request) {
	availableFonts, err := database.GetCustomFonts()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(convertWebfontListToApiFontList(availableFonts))
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func getFontByName(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	font, err := database.GetFont(fontName)

	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(convertWebfontToApiFont(font))
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
	}
}

func createFont(w http.ResponseWriter, r *http.Request) {
	body := new(addFontData)
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	font := database.Webfont{
		Name:        body.Name,
		Fonts:       map[string]database.File{},
		Description: body.Description,
		Designers:   make([]database.Designer, 0),
		License:     body.License,
		Category:    body.Category,
		GoogleFont:  false,
	}

	err = database.CreateFont(&font)
	if err != nil {
		http.Error(w, "Failed to create font", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(font)
	if err != nil {
		http.Error(w, "Failed to encode body", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateFont(w http.ResponseWriter, r *http.Request) {
	fontName := mux.Vars(r)["fontName"]
	font, err := database.GetFont(fontName)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	body := new(updateFontData)

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	font.Description = body.Description
	font.License = body.License
	font.Category = body.Category

	err = database.UpdateFont(font)
	if err != nil {
		http.Error(w, "Failed to update font", http.StatusInternalServerError)
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
		http.Error(w, "Failed to delete font", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func syncFonts(w http.ResponseWriter, _ *http.Request) {
	err := fontsync.Sync()
	if err != nil {
		http.Error(w, "Failed to sync fonts", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
