package admin

import (
	"io/ioutil"
	"jinya-fonts/config"
	"net/http"
	"time"
)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./admin/templates/login.gohtml")
		if err != nil {
			http.NotFound(w, r)
			return
		}

		w.Write(data)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		password := r.FormValue("password")
		if password == config.LoadedConfiguration.AdminPassword {
			cookie := &http.Cookie{
				Name:     "auth",
				Value:    password,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 24),
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			}

			remember := r.FormValue("remember") == "on"
			if remember {
				cookie.Expires = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/admin", http.StatusFound)
		} else {
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
