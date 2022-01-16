package admin

import (
	"jinya-fonts/config"
	"net/http"
	"os"
)

func DeleteFont(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("name")
	if r.Method == http.MethodGet {
		err := RenderAdmin(w, "deleteFont", struct {
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
			RenderAdmin(w, "deleteFont", struct {
				Name    string
				Referer string
				Message string
			}{fontName, referer, "Failed to delete font directory"})
		}
		err = os.Remove(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
		if err != nil {
			RenderAdmin(w, "deleteFont", struct {
				Name    string
				Referer string
				Message string
			}{fontName, referer, "Failed to delete font config file"})
		}

		http.Redirect(w, r, referer, http.StatusFound)
	}
}
