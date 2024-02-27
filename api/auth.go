package api

import (
	"jinya-fonts/config"
	"net/http"
)

func checkAuthorizationHeader(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		password := r.Header.Get("Authorization")
		if password != config.LoadedConfiguration.AdminPassword {
			w.WriteHeader(http.StatusUnauthorized)
		} else {
			handler(w, r)
		}
	}
}
