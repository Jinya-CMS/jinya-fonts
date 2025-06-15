package main

import (
	"embed"
	"github.com/bamzi/jobrunner"
	"github.com/gorilla/mux"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
	"html/template"
	"jinya-fonts/config"
	"jinya-fonts/database"
	"jinya-fonts/fontsync"
	"jinya-fonts/routes/admin"
	"jinya-fonts/routes/api"
	"jinya-fonts/routes/fonts"
	"jinya-fonts/routes/frontend"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
)

var (
	//go:embed openapi
	openapi embed.FS
	//go:embed openapi/admin
	adminOpenapi embed.FS
	//go:embed static
	static embed.FS
)

type SpaHandler struct {
	embedFS      embed.FS
	indexPath    string
	fsPrefixPath string
	templated    bool
	templateData any
}

func (handler SpaHandler) serveTemplated(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFS(handler.embedFS, handler.indexPath)
	if err != nil {
		http.Error(w, "Failed to get admin page", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, handler.templateData)
	if err != nil {
		http.Error(w, "Failed to get admin page", http.StatusInternalServerError)
		return
	}
}

func (handler SpaHandler) servePlain(w http.ResponseWriter, r *http.Request) {
	http.ServeFileFS(w, r, handler.embedFS, handler.indexPath)
}

func (handler SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullPath := strings.TrimPrefix(path.Join(handler.fsPrefixPath, r.URL.Path), "/")
	file, err := handler.embedFS.Open(fullPath)
	if err != nil {
		if handler.templated {
			handler.serveTemplated(w, r)
		} else {
			handler.servePlain(w, r)
		}
		return
	}

	if fi, err := file.Stat(); err != nil || fi.IsDir() {
		if handler.templated {
			handler.serveTemplated(w, r)
		} else {
			handler.servePlain(w, r)
		}
		return
	}

	http.ServeFileFS(w, r, handler.embedFS, fullPath)
}

func main() {
	err := config.LoadConfiguration()
	if err != nil {
		log.Fatalf("Failed with err %v", err)
	}

	database.SetupDatabase()

	settings, err := database.CreateSettingsIfNotExists()
	if err != nil {
		log.Fatalf("Failed with err %v", err)
	}

	jobrunner.Start()

	if settings.SyncEnabled {
		err = jobrunner.Schedule(settings.SyncInterval, fontsync.SyncJob{})
		if err != nil {
			log.Printf("Failed to schedule sync job %s", err.Error())
		}
	}

	if slices.Contains(os.Args, "sync") {
		err = fontsync.Sync()
		if err != nil {
			log.Fatalf("Failed with err %v", err)
		}
	}

	if slices.Contains(os.Args, "serve") {
		router := mux.NewRouter()

		fonts.SetupFontsRouter(router)
		admin.SetupRouter(router)
		api.SetupApiRouter(router)

		router.PathPrefix("/openapi/admin").Handler(SpaHandler{
			embedFS:      adminOpenapi,
			indexPath:    "openapi/admin/index.html",
			fsPrefixPath: "",
			templated:    false,
		})
		router.PathPrefix("/openapi").Handler(SpaHandler{
			embedFS:      openapi,
			indexPath:    "openapi/index.html",
			fsPrefixPath: "",
			templated:    false,
		})

		if config.IsDev() {
			router.PathPrefix("/static/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Service-Worker-Allowed", "/admin")
				http.FileServerFS(os.DirFS(".")).ServeHTTP(w, r)
			})
		} else {
			router.PathPrefix("/static/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Service-Worker-Allowed", "/admin")
				http.FileServerFS(static).ServeHTTP(w, r)
			})
		}

		if config.LoadedConfiguration.ServeWebsite {
			frontend.SetupRouter(router)
		}

		log.Println("Serving at localhost:8090...")
		err = http.ListenAndServe(":8090", router)
		if err != nil {
			log.Fatalf("Failed with err %v", err)
		}
	}
}
