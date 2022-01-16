package admin

import (
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"jinya-fonts/utils"
	"net/http"
	"strings"
)

func checkAuthCookie(r *http.Request) bool {
	authCookie, err := r.Cookie("auth")

	return err != nil || authCookie.Value != config.LoadedConfiguration.AdminPassword
}

func AllFonts(w http.ResponseWriter, r *http.Request) {
	if checkAuthCookie(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type fontData struct {
		Name         string
		NumberStyles int
		License      string
		Category     string
		Author       string
		GoogleFont   bool
	}

	fonts, err := utils.LoadFonts()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var data []fontData
	for _, font := range fonts {
		designers := grabDesigners(font)

		data = append(data, fontData{
			Name:         font.Name,
			NumberStyles: len(font.Fonts),
			License:      font.License,
			Category:     font.Category,
			Author:       strings.Join(designers, ","),
			GoogleFont:   font.GoogleFont,
		})
	}

	err = RenderAdmin(w, "allFonts", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func CustomFonts(w http.ResponseWriter, r *http.Request) {
	if checkAuthCookie(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type fontData struct {
		Name         string
		NumberStyles int
		License      string
		Category     string
		Author       string
	}

	fonts, err := utils.LoadFonts()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var data []fontData
	for _, font := range fonts {
		if font.GoogleFont {
			continue
		}

		designers := grabDesigners(font)

		data = append(data, fontData{
			Name:         font.Name,
			NumberStyles: len(font.Fonts),
			License:      font.License,
			Category:     font.Category,
			Author:       strings.Join(designers, ","),
		})
	}

	err = RenderAdmin(w, "customFonts", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SyncedFonts(w http.ResponseWriter, r *http.Request) {
	if checkAuthCookie(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type fontData struct {
		Name         string
		NumberStyles int
		License      string
		Category     string
		Author       string
	}

	fonts, err := utils.LoadFonts()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var data []fontData
	for _, font := range fonts {
		if !font.GoogleFont {
			continue
		}

		designers := grabDesigners(font)

		data = append(data, fontData{
			Name:         font.Name,
			NumberStyles: len(font.Fonts),
			License:      font.License,
			Category:     font.Category,
			Author:       strings.Join(designers, ","),
		})
	}

	err = RenderAdmin(w, "syncedFonts", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func grabDesigners(font meta.FontFile) []string {
	var designers []string
	for _, designer := range font.Designers {
		designers = append(designers, designer.Name)
	}
	return designers
}
