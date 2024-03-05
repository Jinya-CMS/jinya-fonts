package api

import (
	"cmp"
	"encoding/json"
	"github.com/gorilla/mux"
	"jinya-fonts/database"
	"net/http"
	"slices"
)

type apiFont struct {
	Name        string              `json:"name" bson:"name,omitempty"`
	Description string              `json:"description" bson:"description,omitempty"`
	Designers   []database.Designer `json:"designers" bson:"designers,omitempty"`
	License     string              `json:"license" bson:"license,omitempty"`
	Category    string              `json:"category" bson:"category,omitempty"`
	Fonts       []database.File     `json:"fonts"`
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

func getFonts(w http.ResponseWriter, r *http.Request) {
	availableFonts, err := database.GetAllFonts()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = json.NewEncoder(w).Encode(convertWebfontListToApiFontList(availableFonts))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
