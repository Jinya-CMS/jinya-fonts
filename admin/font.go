package admin

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"jinya-fonts/utils"
	"net/http"
	"os"
	"strings"
)

func AddFont(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := RenderAdmin(w, "font/add", struct {
			Message  string
			Name     string
			License  string
			Category string
			Referer  string
		}{"", "", "", "", r.Referer()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			RenderAdmin(w, "font/add", struct {
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
			RenderAdmin(w, "font/add", struct {
				Message  string
				Name     string
				License  string
				Category string
				Referer  string
			}{"Please provide a name", name, license, category, referer})
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
				Message  string
				Name     string
				License  string
				Category string
				Referer  string
			}{"A font with the given name exists already", name, license, category, referer})
			return
		}

		fontFile := meta.FontFile{
			Name:        name,
			Fonts:       []meta.FontFileMeta{},
			Description: description,
			Designers:   []meta.FontDesigner{},
			License:     license,
			Category:    category,
			GoogleFont:  false,
		}

		data, err := yaml.Marshal(fontFile)
		if err != nil {
			RenderAdmin(w, "font/add", struct {
				Message  string
				Name     string
				License  string
				Category string
				Referer  string
			}{"Failed to convert font metadata", name, license, category, referer})
			return
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+name+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "font/add", struct {
				Message  string
				Name     string
				License  string
				Category string
				Referer  string
			}{"Failed to save font metadata", name, license, category, referer})
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
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + r.URL.Query().Get("name") + ".yaml")
	if err != nil {
		RenderAdmin(w, "font/edit", struct {
			Message  string
			Name     string
			License  string
			Category string
			Referer  string
		}{"Font does not exist", "", "", "", r.Referer()})
		return
	}

	var font meta.FontFile
	err = yaml.Unmarshal(data, &font)
	if err != nil {
		RenderAdmin(w, "font/edit", struct {
			Message  string
			Name     string
			License  string
			Category string
			Referer  string
		}{"Font does not exist", "", "", "", r.Referer()})
		return
	}

	if r.Method == http.MethodGet {
		err = RenderAdmin(w, "font/edit", struct {
			Message  string
			Name     string
			License  string
			Category string
			Referer  string
		}{"", font.Name, font.License, font.Category, r.Referer()})
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
				Message  string
				Name     string
				License  string
				Category string
				Referer  string
			}{"Please provide a name", name, license, category, referer})
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
				Message  string
				Name     string
				License  string
				Category string
				Referer  string
			}{"Failed to convert font metadata", name, license, category, referer})
			return
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+name+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "font/edit", struct {
				Message  string
				Name     string
				License  string
				Category string
				Referer  string
			}{"Failed to save font metadata", name, license, category, referer})
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

	err = RenderAdmin(w, "font/synced", data)
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
