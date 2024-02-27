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
	"mime"
	"net/http"
	"os"
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

func getAngularStatic(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/v3/")
	var data []byte
	var err error

	if strings.HasSuffix(path, ".js") || strings.HasSuffix(path, ".css") || strings.HasSuffix(path, ".png") || strings.HasSuffix(path, ".ico") {
		path = strings.TrimPrefix(path, "static/")
		data, err = angular.ReadFile("angular/frontend/dist/browser/" + path)
		if err != nil {
			log.Printf("page %s not found in pages cache...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		data, err = angular.ReadFile("angular/frontend/dist/browser/index.html")
		if err != nil {
			log.Printf("page %s not found in pages cache...", r.RequestURI)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	lastIdx := strings.LastIndex(path, ".")
	extension := "text/html"
	if lastIdx > 0 {
		extension = path[lastIdx:len(path)]
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(extension))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getOpenApiYaml(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var err error

	data, err = openapi.ReadFile("openapi/v3/jinya-fonts.yaml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func getOpenApi(w http.ResponseWriter, r *http.Request) {
	var data []byte
	var err error

	data, err = openapi.ReadFile("openapi/v3/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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

			router.PathPrefix("/v3/").HandlerFunc(getAngularStatic)

			router.HandleFunc("/openapi.yaml", getOpenApiYaml)
			router.HandleFunc("/openapi", getOpenApi)
			router.PathPrefix("/openapi/static/").Handler(http.FileServer(http.FS(openapi)))

			router.PathPrefix("/static/").Handler(http.FileServer(http.FS(static)))
			router.PathPrefix("/").HandlerFunc(getWebApp)
		}

		http.Handle("/", router)

		log.Println("Serving at localhost:8090...")
		err = http.ListenAndServe(":8090", nil)
		if err != nil {
			panic(err)
		}
	}
}
