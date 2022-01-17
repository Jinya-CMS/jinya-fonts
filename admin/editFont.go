package admin

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
)

func EditFont(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + r.URL.Query().Get("name") + ".yaml")
	if err != nil {
		RenderAdmin(w, "editFont", struct {
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
		RenderAdmin(w, "editFont", struct {
			Message  string
			Name     string
			License  string
			Category string
			Referer  string
		}{"Font does not exist", "", "", "", r.Referer()})
		return
	}

	if r.Method == http.MethodGet {
		err = RenderAdmin(w, "editFont", struct {
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
			RenderAdmin(w, "editFont", struct {
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
			RenderAdmin(w, "editFont", struct {
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
			RenderAdmin(w, "editFont", struct {
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
			RenderAdmin(w, "editFont", struct {
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
