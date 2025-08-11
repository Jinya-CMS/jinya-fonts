package fonts

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupFontsRouter(router *mux.Router) {
	router.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts", fileHandler{}))
	router.Methods("GET").Path("/css2").HandlerFunc(getCss2)
	router.HandleFunc("/download", downloadFont)
}
