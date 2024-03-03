package main

import (
	"embed"
	"github.com/gorilla/mux"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
	"html/template"
	admin "jinya-fonts/admin/api"
	"jinya-fonts/api"
	"jinya-fonts/config"
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
}

func (handler SpaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fullPath := strings.TrimPrefix(path.Join(handler.fsPrefixPath, r.URL.Path), "/")
	file, err := handler.embedFS.Open(fullPath)
	if err != nil {
		http.ServeFileFS(w, r, handler.embedFS, handler.indexPath)
		return
	}

	if fi, err := file.Stat(); err != nil || fi.IsDir() {
		http.ServeFileFS(w, r, handler.embedFS, handler.indexPath)
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

		router.PathPrefix("/openapi/admin").Handler(SpaHandler{
			embedFS:      adminOpenapi,
			indexPath:    "openapi/admin/index.html",
			fsPrefixPath: "",
		})
		router.PathPrefix("/openapi").Handler(SpaHandler{
			embedFS:      openapi,
			indexPath:    "openapi/index.html",
			fsPrefixPath: "",
		})

		if config.LoadedConfiguration.ServeWebsite {
			api.SetupApiRouter(router)

			router.PathPrefix("/v3").Handler(http.StripPrefix("/v3", SpaHandler{
				embedFS:      angular,
				indexPath:    "angular/frontend/dist/browser/index.html",
				fsPrefixPath: "angular/frontend/dist/browser",
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
