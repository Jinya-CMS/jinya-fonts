package api

import "github.com/gorilla/mux"

func SetupApiRouter(router *mux.Router) {
	router.Methods("GET").Path("/api/font").HandlerFunc(getAllFonts)
	router.Methods("POST").Path("/api/font").HandlerFunc(checkAuthorizationHeader(createFont))
	router.Methods("GET").Path("/api/font/{fontName}").HandlerFunc(getFontByName)
	router.Methods("PUT").Path("/api/font/{fontName}").HandlerFunc(checkAuthorizationHeader(updateFont))
	router.Methods("DELETE").Path("/api/font/{fontName}").HandlerFunc(checkAuthorizationHeader(deleteFont))

	router.Methods("GET").Path("/api/font/{fontName}/file").HandlerFunc(getFontFiles)
	router.Methods("POST").Path("/api/font/{fontName}/file/{fontSubset}/{fontWeight}/{fontStyle}").HandlerFunc(checkAuthorizationHeader(createFontFile))
	router.Methods("PUT").Path("/api/font/{fontName}/file/{fontSubset}/{fontWeight}/{fontStyle}").HandlerFunc(checkAuthorizationHeader(updateFontFile))
	router.Methods("DELETE").Path("/api/font/{fontName}/file/{fontSubset}/{fontWeight}/{fontStyle}").HandlerFunc(checkAuthorizationHeader(deleteFontFile))

	router.Methods("GET").Path("/api/font/{fontName}/designer").HandlerFunc(getFontDesigners)
	router.Methods("POST").Path("/api/font/{fontName}/designer").HandlerFunc(checkAuthorizationHeader(createFontDesigner))
	router.Methods("DELETE").Path("/api/font/{fontName}/designer/{designerName}").HandlerFunc(checkAuthorizationHeader(deleteFontDesigner))
}
