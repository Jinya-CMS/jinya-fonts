package api

import (
	"github.com/gorilla/mux"
	"jinya-fonts/routes/admin/api"
	"net/http"
)

func SetupApiRouter(router *mux.Router) {
	router.Methods("GET").Path("/api/font").HandlerFunc(api.GetAllFonts)
	router.Methods("GET").Path("/api/font/{fontName}").HandlerFunc(api.GetFontByName)

	router.Methods("GET").Path("/api/font/{fontName}/file").HandlerFunc(api.GetFontFiles)

	router.Methods("GET").Path("/api/font/{fontName}/designer").HandlerFunc(api.GetFontDesigners)

	router.Methods("GET").Path("/healthz").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}
