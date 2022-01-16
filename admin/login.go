package admin

import (
	"io/ioutil"
	"jinya-fonts/config"
	"net/http"
	"strings"
	"time"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := strings.TrimPrefix(r.URL.Path, "/login")
		if path == "/" || path == "" {
			path = "login.html"
		}

		data, err := ioutil.ReadFile("./admin/" + path)
		if err != nil {
			path = "login.html"
			data, err = ioutil.ReadFile("./admin/" + path)

			if err != nil {
				http.NotFound(w, r)
				return
			}
		}

		if strings.HasSuffix(path, "css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, "html") {
			w.Header().Set("Content-Type", "text/html")
		} else if strings.HasSuffix(path, "js") {
			w.Header().Set("Content-Type", "application/javascript")
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
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
