package api

import "github.com/gorilla/mux"

func SetupApiRouter(router *mux.Router) {
	router.Methods("GET").Path("/api/font").HandlerFunc(getFonts)
	router.Methods("GET").Path("/api/font/{fontName}").HandlerFunc(getFontByName)

	router.Methods("GET").Path("/api/font/{fontName}/file").HandlerFunc(getFontFiles)

	router.Methods("GET").Path("/api/font/{fontName}/designer").HandlerFunc(getFontDesigners)
}
