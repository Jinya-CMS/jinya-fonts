package api

import (
	"jinya-fonts/routes/admin/api"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupApiRouter(router *mux.Router) {
	router.
		Methods(http.MethodGet).
		Path("/api/font").
		HandlerFunc(api.GetAllFonts)
	router.
		Methods(http.MethodGet).
		Path("/api/font/{fontName}").
		HandlerFunc(api.GetFontByName)

	router.
		Methods(http.MethodGet).
		Path("/api/font/{fontName}/file").
		HandlerFunc(api.GetFontFiles)

	router.
		Methods(http.MethodGet).
		Path("/api/font/{fontName}/designer").
		HandlerFunc(api.GetFontDesigners)

	router.
		Methods(http.MethodGet).
		Path("/healthz").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
}
