package main

import (
	"embed"
	"flag"
	"html/template"
	"jinya-fonts/admin"
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
	//go:embed admin/static
	adminStatic embed.FS
	//go:embed static
	static embed.FS
	pages  = map[string]string{
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
		http.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir(configuration.FontFileFolder))))
		http.HandleFunc("/css2", http2.GetCss2)
		if configuration.AdminPassword != "" {
			http.HandleFunc("/login", admin.Login)
			http.HandleFunc("/login/", admin.Login)
			http.HandleFunc("/logout", admin.Logout)
			http.HandleFunc("/logout/", admin.Logout)

			http.HandleFunc("/admin", admin.CheckAuthCookie(admin.AllFonts))
			http.HandleFunc("/admin/add", admin.CheckAuthCookie(admin.AddFont))
			http.HandleFunc("/admin/edit", admin.CheckAuthCookie(admin.EditFont))
			http.HandleFunc("/admin/delete", admin.CheckAuthCookie(admin.DeleteFont))
			http.HandleFunc("/admin/sync", admin.CheckAuthCookie(admin.TriggerSync))
			http.HandleFunc("/admin/synced", admin.CheckAuthCookie(admin.SyncedFonts))
			http.HandleFunc("/admin/custom", admin.CheckAuthCookie(admin.CustomFonts))

			http.HandleFunc("/admin/designers", admin.CheckAuthCookie(admin.DesignersIndex))
			http.HandleFunc("/admin/designers/delete", admin.CheckAuthCookie(admin.DeleteDesigner))
			http.HandleFunc("/admin/designers/add", admin.CheckAuthCookie(admin.AddDesigner))
			http.HandleFunc("/admin/designers/edit", admin.CheckAuthCookie(admin.EditDesigner))

			http.HandleFunc("/admin/files", admin.CheckAuthCookie(admin.FilesIndex))
			http.HandleFunc("/admin/files/delete", admin.CheckAuthCookie(admin.DeleteFile))
			http.HandleFunc("/admin/files/add", admin.CheckAuthCookie(admin.AddFile))
			http.HandleFunc("/admin/files/edit", admin.CheckAuthCookie(admin.EditFile))

			http.Handle("/admin/static/", http.FileServer(http.FS(adminStatic)))
		}
		if configuration.ServeWebsite {
			http.HandleFunc("/api/font", http2.GetFontMeta)
			http.HandleFunc("/download", http2.DownloadFont)

			http.HandleFunc("/v3/", getAngularStatic)

			http.Handle("/static/", http.FileServer(http.FS(static)))
			http.HandleFunc("/", getWebApp)
		}

		log.Println("Serving at localhost:8090...")
		err = http.ListenAndServe(":8090", nil)
		if err != nil {
			panic(err)
		}
	}
}
