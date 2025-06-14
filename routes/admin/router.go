package admin

import (
	"github.com/gorilla/mux"
	"jinya-fonts/routes/admin/api"
	"jinya-fonts/routes/admin/web"
)

func SetupRouter(router *mux.Router) {
	api.SetupAdminApiRouter(router)
	web.SetupRouter(router)
}
