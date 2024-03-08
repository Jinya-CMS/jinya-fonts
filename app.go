package main

import (
	"embed"
	"errors"
	"github.com/bamzi/jobrunner"
	"github.com/gorilla/mux"
	_ "github.com/sakirsensoy/genv/dotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/language"
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
)

type LanguageHandler struct {
	defaultLang     language.Tag
	langPathMapping map[language.Tag]string
}

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

func (handler LanguageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	oldPath := "/" + r.URL.Path + "?" + r.URL.RawQuery
	acceptLanguage, _, err := language.ParseAcceptLanguage(r.Header.Get("Accept-Language"))
	if err != nil {
		http.Redirect(w, r, handler.langPathMapping[handler.defaultLang]+oldPath, http.StatusFound)
		return
	}

	localMap := map[string]string{}
	for tag, p := range handler.langPathMapping {
		b, _ := tag.Base()
		localMap[b.ISO3()] = p
	}

	for _, lang := range acceptLanguage {
		b, _ := lang.Base()
		if p, exists := localMap[b.ISO3()]; exists {
			http.Redirect(w, r, p+oldPath, http.StatusFound)
			return
		}
	}

	http.Redirect(w, r, handler.langPathMapping[handler.defaultLang]+oldPath, http.StatusFound)
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
	} else if err != nil {
		panic(err)
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
		router.PathPrefix("/admin/de").Handler(http.StripPrefix("/admin/de", SpaHandler{
			embedFS:      angularAdmin,
			indexPath:    "admin/web/dist/browser/de/index.html",
			fsPrefixPath: "admin/web/dist/browser/de",
			templated:    true,
			templateData: config.LoadedConfiguration,
		}))
		router.PathPrefix("/admin/en").Handler(http.StripPrefix("/admin/en", SpaHandler{
			embedFS:      angularAdmin,
			indexPath:    "admin/web/dist/browser/en/index.html",
			fsPrefixPath: "admin/web/dist/browser/en",
			templated:    true,
			templateData: config.LoadedConfiguration,
		}))
		router.PathPrefix("/admin").Handler(http.StripPrefix("/admin", LanguageHandler{
			defaultLang: language.English,
			langPathMapping: map[language.Tag]string{
				language.English: "/admin/en",
				language.German:  "/admin/de",
			},
		}))

		if config.LoadedConfiguration.ServeWebsite {
			router.PathPrefix("/static/").Handler(http.FileServerFS(static))
			router.PathPrefix("/de").Handler(http.StripPrefix("/de", SpaHandler{
				embedFS:      angular,
				indexPath:    "angular/frontend/dist/browser/de/index.html",
				fsPrefixPath: "angular/frontend/dist/browser/de",
				templated:    false,
			}))
			router.PathPrefix("/en").Handler(http.StripPrefix("/en", SpaHandler{
				embedFS:      angular,
				indexPath:    "angular/frontend/dist/browser/en/index.html",
				fsPrefixPath: "angular/frontend/dist/browser/en",
				templated:    false,
			}))
			router.PathPrefix("/").Handler(LanguageHandler{
				defaultLang: language.English,
				langPathMapping: map[language.Tag]string{
					language.English: "/en",
					language.German:  "/de",
				},
			})
		}

		log.Println("Serving at localhost:8090...")
		err = http.ListenAndServe(":8090", router)
		if err != nil {
			panic(err)
		}
	}
}
