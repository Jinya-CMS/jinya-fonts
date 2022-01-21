package admin

import (
	"jinya-fonts/config"
	"net/http"
)

func CheckAuthCookie(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		authCookie, err := r.Cookie("auth")
		if err != nil || authCookie.Value != config.LoadedConfiguration.AdminPassword {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			handler(w, r)
		}
	}
}
