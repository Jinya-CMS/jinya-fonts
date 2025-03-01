package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func SetupApiRouter(router *mux.Router) {
	router.Methods("GET").Path("/api/font").HandlerFunc(getFonts)
	router.Methods("GET").Path("/api/font/{fontName}").HandlerFunc(getFontByName)

	router.Methods("GET").Path("/api/font/{fontName}/file").HandlerFunc(getFontFiles)

	router.Methods("GET").Path("/api/font/{fontName}/designer").HandlerFunc(getFontDesigners)

	router.Methods("GET").Path("/healthz").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
