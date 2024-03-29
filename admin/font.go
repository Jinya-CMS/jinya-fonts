package admin

import (
	"bytes"
	"gopkg.in/yaml.v3"
	"jinya-fonts/config"
	"jinya-fonts/database"
	"jinya-fonts/fontsync"
	"log"
	"net/http"
	"os"
	"strings"
)

func AddFont(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := RenderAdmin(w, "font/add", struct {
			Message     string
			Name        string
			License     string
			Category    string
			Description string
			Referer     string
		}{"", "", "", "", "", r.Referer()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			RenderAdmin(w, "font/add", struct {
				Message     string
				Name        string
				License     string
				Category    string
				Description string
				Referer     string
			}{"Failed to parse input", "", "", "", "", r.Referer()})
			return
		}

		name := r.FormValue("name")
		license := r.FormValue("license")
		category := r.FormValue("category")
		description := r.FormValue("description")
		referer := r.FormValue("referer")

		if name == "" {
			RenderAdmin(w, "font/add", struct {
				Message     string
				Name        string
				License     string
				Category    string
				Description string
				Referer     string
			}{"Please provide a name", name, license, category, description, referer})
			return
		}
		if license == "" {
			license = "All rights reserved"
		}
		if category == "" {
			category = "Sans Serif"
		}

		if _, err := os.Stat(config.LoadedConfiguration.FontFileFolder + "/" + name + ".yaml"); !os.IsNotExist(err) {
			RenderAdmin(w, "font/add", struct {
				Message     string
				Name        string
				License     string
				Category    string
				Description string
				Referer     string
			}{"A font with the given name exists already", name, license, category, description, referer})
			return
		}

		fontFile := database.Webfont{
			Name:        name,
			Fonts:       []database.Metadata{},
			Description: description,
			Designers:   []database.Designer{},
			License:     license,
			Category:    category,
			GoogleFont:  false,
		}

		data, err := yaml.Marshal(fontFile)
		if err != nil {
			RenderAdmin(w, "font/add", struct {
				Message     string
				Name        string
				License     string
				Category    string
				Description string
				Referer     string
			}{"Failed to convert font metadata", name, license, category, description, referer})
			return
		}

		err = os.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+name+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "font/add", struct {
				Message     string
				Name        string
				License     string
				Category    string
				Description string
				Referer     string
			}{"Failed to save font metadata", name, license, category, description, referer})
			return
		}

		http.Redirect(w, r, referer, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func DeleteFont(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("name")
	if r.Method == http.MethodGet {
		err := RenderAdmin(w, "font/delete", struct {
			Name    string
			Referer string
			Message string
		}{fontName, r.Referer(), ""})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		referer := r.FormValue("referer")
		if r.FormValue("delete") != "true" {
			http.Redirect(w, r, referer, http.StatusFound)
			return
		}

		err = os.RemoveAll(config.LoadedConfiguration.FontFileFolder + "/" + fontName)
		if err != nil {
			RenderAdmin(w, "font/delete", struct {
				Name    string
				Referer string
				Message string
			}{fontName, referer, "Failed to delete font directory"})
		}
		err = os.Remove(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
		if err != nil {
			RenderAdmin(w, "font/delete", struct {
				Name    string
				Referer string
				Message string
			}{fontName, referer, "Failed to delete font config file"})
		}

		http.Redirect(w, r, referer, http.StatusFound)
	}
}

func EditFont(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + r.URL.Query().Get("name") + ".yaml")
	if err != nil {
		RenderAdmin(w, "font/edit", struct {
			Message     string
			Name        string
			License     string
			Category    string
			Description string
			Referer     string
		}{"Font does not exist", "", "", "", "", r.Referer()})
		return
	}

	var font database.Webfont
	err = yaml.Unmarshal(data, &font)
	if err != nil {
		RenderAdmin(w, "font/edit", struct {
			Message     string
			Name        string
			License     string
			Category    string
			Description string
			Referer     string
		}{"Font does not exist", "", "", "", "", r.Referer()})
		return
	}

	if r.Method == http.MethodGet {
		err = RenderAdmin(w, "font/edit", struct {
			Message     string
			Name        string
			License     string
			Category    string
			Description string
			Referer     string
		}{"", font.Name, font.License, font.Category, font.Description, r.Referer()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			RenderAdmin(w, "font/edit", struct {
				Message  string
				Name     string
				License  string
				Category string
				Referer  string
			}{"Failed to parse input", "", "", "", r.Referer()})
			return
		}

		name := r.FormValue("name")
		license := r.FormValue("license")
		category := r.FormValue("category")
		description := r.FormValue("description")
		referer := r.FormValue("referer")

		if name == "" {
			RenderAdmin(w, "font/edit", struct {
				Message     string
				Name        string
				License     string
				Category    string
				Description string
				Referer     string
			}{"Please provide a name", name, license, category, description, referer})
			return
		}
		if license == "" {
			license = "All rights reserved"
		}
		if category == "" {
			category = "Sans Serif"
		}

		font.Name = name
		font.License = license
		font.Category = category
		font.Description = description

		data, err := yaml.Marshal(font)
		if err != nil {
			RenderAdmin(w, "font/edit", struct {
				Message     string
				Name        string
				License     string
				Category    string
				Description string
				Referer     string
			}{"Failed to convert font metadata", name, license, category, description, referer})
			return
		}

		err = os.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+name+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "font/edit", struct {
				Message     string
				Name        string
				License     string
				Category    string
				Description string
				Referer     string
			}{"Failed to save font metadata", name, license, category, description, referer})
			return
		}

		http.Redirect(w, r, referer, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AllFonts(w http.ResponseWriter, r *http.Request) {
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

	fonts, err := database.GetAllFonts()
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
			Author:       strings.Join(designers, ", "),
			GoogleFont:   font.GoogleFont,
		})
	}

	err = RenderAdmin(w, "font/all", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func CustomFonts(w http.ResponseWriter, r *http.Request) {
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

	fonts, err := database.GetAllFonts()
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
			Author:       strings.Join(designers, ", "),
		})
	}

	err = RenderAdmin(w, "font/custom", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func SyncedFonts(w http.ResponseWriter, r *http.Request) {
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

	fonts, err := database.GetAllFonts()
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
			Author:       strings.Join(designers, ", "),
		})
	}

	err = RenderAdmin(w, "font/synced", data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func grabDesigners(font database.Webfont) []string {
	var designers []string
	for _, designer := range font.Designers {
		designers = append(designers, designer.Name)
	}
	return designers
}

func TriggerSync(w http.ResponseWriter, r *http.Request) {
	type syncData struct {
		Log string
	}

	if r.Method == http.MethodPost {
		backupWriter := log.Writer()
		buffer := bytes.Buffer{}
		log.SetOutput(&buffer)
		_, err := config.LoadConfiguration(config.ConfigurationPath)
		if err != nil {
			log.Println("Failed to reparse config")
		}
		_ = fontsync.Sync(config.LoadedConfiguration)

		log.SetOutput(backupWriter)

		data := syncData{Log: buffer.String()}
		err = RenderAdmin(w, "font/triggerSync", data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := RenderAdmin(w, "font/triggerSync", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
