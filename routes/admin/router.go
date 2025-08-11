package admin

import (
	"jinya-fonts/routes/admin/api"
	"jinya-fonts/routes/admin/web"

	"github.com/gorilla/mux"
)

func SetupRouter(router *mux.Router) {
	api.SetupAdminApiRouter(router)
	web.SetupRouter(router)
}
