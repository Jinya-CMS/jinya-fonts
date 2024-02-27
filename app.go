package main

import (
	"embed"
	"flag"
	"github.com/gorilla/mux"
	"html/template"
	"jinya-fonts/admin"
	"jinya-fonts/api"
	"jinya-fonts/config"
	"jinya-fonts/fontsync"
	http2 "jinya-fonts/http"
	"jinya-fonts/utils"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	//go:embed frontend
	frontend embed.FS
	//go:embed angular/frontend/dist/browser
	angular embed.FS
	//go:embed openapi
	openapi embed.FS
	//go:embed static
	static embed.FS
	//go:embed admin/static
	adminStatic embed.FS
	pages       = map[string]string{
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
	if _, err := handler.embedFS.Open(fullPath); err != nil {
		http.ServeFileFS(w, r, handler.embedFS, handler.indexPath)
		return
	}

	http.FileServerFS(handler.embedFS).ServeHTTP(w, r)
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
	configFileFlag := flag.String("config-file", "./config.yaml", "The config file, check the sample for the structure")
	flag.Parse()

	configuration, err := config.LoadConfiguration(*configFileFlag)
	if err != nil {
		panic(err)
	}

	if utils.ContainsString(os.Args, "sync") {
		err = fontsync.Sync(configuration)
		if err != nil {
			panic(err)
		}
	}

	if utils.ContainsString(os.Args, "serve") {
		router := mux.NewRouter()

		router.PathPrefix("/fonts/").Handler(http.StripPrefix("/fonts", http.FileServer(http.Dir(configuration.FontFileFolder))))
		router.HandleFunc("/css2", http2.GetCss2)
		if configuration.AdminPassword != "" {
			router.HandleFunc("/login", admin.Login)
			router.HandleFunc("/logout", admin.Logout)

			router.HandleFunc("/admin", admin.CheckAuthCookie(admin.AllFonts))
			router.HandleFunc("/admin/add", admin.CheckAuthCookie(admin.AddFont))
			router.HandleFunc("/admin/edit", admin.CheckAuthCookie(admin.EditFont))
			router.HandleFunc("/admin/delete", admin.CheckAuthCookie(admin.DeleteFont))
			router.HandleFunc("/admin/sync", admin.CheckAuthCookie(admin.TriggerSync))
			router.HandleFunc("/admin/synced", admin.CheckAuthCookie(admin.SyncedFonts))
			router.HandleFunc("/admin/custom", admin.CheckAuthCookie(admin.CustomFonts))

			router.HandleFunc("/admin/designers", admin.CheckAuthCookie(admin.DesignersIndex))
			router.HandleFunc("/admin/designers/delete", admin.CheckAuthCookie(admin.DeleteDesigner))
			router.HandleFunc("/admin/designers/add", admin.CheckAuthCookie(admin.AddDesigner))
			router.HandleFunc("/admin/designers/edit", admin.CheckAuthCookie(admin.EditDesigner))

			router.HandleFunc("/admin/files", admin.CheckAuthCookie(admin.FilesIndex))
			router.HandleFunc("/admin/files/delete", admin.CheckAuthCookie(admin.DeleteFile))
			router.HandleFunc("/admin/files/add", admin.CheckAuthCookie(admin.AddFile))
			router.HandleFunc("/admin/files/edit", admin.CheckAuthCookie(admin.EditFile))

			router.PathPrefix("/admin/static/").Handler(http.FileServer(http.FS(adminStatic)))
		}
		if configuration.ServeWebsite {
			api.SetupApiRouter(router)

			router.HandleFunc("/download", http2.DownloadFont)

			router.PathPrefix("/v3").Handler(SpaHandler{
				embedFS:      angular,
				indexPath:    "angular/frontend/dist/browser/index.html",
				fsPrefixPath: "angular/frontend/dist/browser",
			})

			router.PathPrefix("/openapi").Handler(SpaHandler{
				embedFS:      openapi,
				indexPath:    "openapi/index.html",
				fsPrefixPath: "",
			})

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
