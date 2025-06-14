package web

import (
	"github.com/gorilla/mux"
)

func SetupRouter(router *mux.Router) {
	router.PathPrefix("/admin/").HandlerFunc(IndexPage)
	router.Path("/admin").HandlerFunc(IndexPage)
}
