package api

import (
	"github.com/gorilla/mux"
)

func SetupAdminApiRouter(router *mux.Router) {
	router.Methods("GET").Path("/api/admin/font/all").HandlerFunc(getAllFonts)
	router.Methods("GET").Path("/api/admin/font/google").HandlerFunc(getGoogleFonts)
	router.Methods("GET").Path("/api/admin/font/custom").HandlerFunc(getCustomFonts)
	router.Methods("POST").Path("/api/admin/font").HandlerFunc(createFont)
	router.Methods("GET").Path("/api/admin/font/{fontName}").HandlerFunc(getFontByName)
	router.Methods("PUT").Path("/api/admin/font/{fontName}").HandlerFunc(updateFont)
	router.Methods("DELETE").Path("/api/admin/font/{fontName}").HandlerFunc(deleteFont)

	router.Methods("GET").Path("/api/admin/font/{fontName}/file").HandlerFunc(getFontFiles)
	router.Methods("POST").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").HandlerFunc(createFontFile)
	router.Methods("PUT").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").HandlerFunc(updateFontFile)
	router.Methods("DELETE").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").HandlerFunc(deleteFontFile)

	router.Methods("GET").Path("/api/admin/font/{fontName}/designer").HandlerFunc(getFontDesigners)
	router.Methods("POST").Path("/api/admin/font/{fontName}/designer").HandlerFunc(createFontDesigner)
	router.Methods("DELETE").Path("/api/admin/font/{fontName}/designer/{designerName}").HandlerFunc(deleteFontDesigner)

	router.Methods("GET").Path("/api/admin/settings").HandlerFunc(getSettings)
	router.Methods("PUT").Path("/api/admin/settings").HandlerFunc(updateSettings)
}
