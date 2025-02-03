package admin

import (
	"github.com/gorilla/mux"
	"jinya-fonts/admin/api"
	"jinya-fonts/admin/web"
)

func SetupRouter(router *mux.Router) {
	api.SetupAdminApiRouter(router)
	web.SetupRouter(router)
}
