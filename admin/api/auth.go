package api

import (
	"net/http"
)

func checkAuthorizationHeader(handler func(w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}
}
