package admin

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
	"os"
)

func AddFont(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		err := RenderAdmin(w, "addFont", struct {
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
			RenderAdmin(w, "addFont", struct {
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
			RenderAdmin(w, "addFont", struct {
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
			RenderAdmin(w, "addFont", struct {
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
			RenderAdmin(w, "addFont", struct {
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
			RenderAdmin(w, "addFont", struct {
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
