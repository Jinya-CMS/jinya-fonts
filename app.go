package main

import (
	"embed"
	"errors"
	"github.com/bamzi/jobrunner"
	"github.com/gorilla/mux"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	admin "jinya-fonts/admin/api"
	"jinya-fonts/api"
	"jinya-fonts/config"
	"jinya-fonts/database"
	"jinya-fonts/fonts"
	"jinya-fonts/fontsync"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"
)

var (
	//go:embed frontend
	frontend embed.FS
	//go:embed angular/frontend/dist/browser
	angular embed.FS
	//go:embed admin/web/dist/browser
	angularAdmin embed.FS
	//go:embed openapi
	openapi embed.FS
	//go:embed openapi/admin
	adminOpenapi embed.FS
	//go:embed static
	static embed.FS
	pages  = map[string]string{
		"/":     "frontend/index.gohtml",
		"/font": "frontend/font.gohtml",
	}
)

type SpaHandler struct {
	embedFS      embed.FS
	indexPath    string
	fsPrefixPath string
	templated    bool
	templateData any
}

func (handler SpaHandler) serveTemplated(w http.ResponseWriter, r *http.Request) {
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

func getWebApp(w http.ResponseWriter, r *http.Request) {
	page, ok := pages[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tpl, err := template.New("layout").ParseFS(frontend, "frontend/layout.gohtml", page)
	if err != nil {
		log.Printf("page %s not found in pages cache...", r.RequestURI)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if err := tpl.Execute(w, nil); err != nil {
		return
	}
}

func main() {
	err := config.LoadConfiguration()
	if err != nil {
		panic(err)
	}

	settings, err := database.GetSettings()
	if errors.Is(err, mongo.ErrNoDocuments) {
		settings = &database.JinyaFontsSettings{
			FilterByName: []string{},
			SyncEnabled:  true,
			SyncInterval: "0 0 1 * *",
		}

		err = database.UpdateSettings(settings)
		if err != nil {
			panic(err)
		}
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
			panic(err)
		}
	}

	if slices.Contains(os.Args, "serve") {
		router := mux.NewRouter()

		fonts.SetupFontsRouter(router)
		admin.SetupAdminApiRouter(router)
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
		router.PathPrefix("/admin").Handler(http.StripPrefix("/admin", SpaHandler{
			embedFS:      angularAdmin,
			indexPath:    "admin/web/dist/browser/index.html",
			fsPrefixPath: "admin/web/dist/browser",
			templated:    true,
			templateData: config.LoadedConfiguration,
		}))

		if config.LoadedConfiguration.ServeWebsite {
			router.PathPrefix("/v3").Handler(http.StripPrefix("/v3", SpaHandler{
				embedFS:      angular,
				indexPath:    "angular/frontend/dist/browser/index.html",
				fsPrefixPath: "angular/frontend/dist/browser",
				templated:    false,
			}))

			router.PathPrefix("/static/").Handler(http.FileServerFS(static))
			router.PathPrefix("/").HandlerFunc(getWebApp)
		}

		log.Println("Serving at localhost:8090...")
		err = http.ListenAndServe(":8090", router)
		if err != nil {
			panic(err)
		}
	}
}
