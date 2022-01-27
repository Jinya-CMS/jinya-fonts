package main

import (
	"flag"
	"jinya-fonts/admin"
	"jinya-fonts/config"
	"jinya-fonts/fontsync"
	http2 "jinya-fonts/http"
	"jinya-fonts/utils"
	"log"
	"net/http"
	"os"
)

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
		http.HandleFunc("/fonts/", http2.GetFont)
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
			http.HandleFunc("/admin/files/delete", admin.CheckAuthCookie(admin.DeleteDesigner))
			http.HandleFunc("/admin/files/add", admin.CheckAuthCookie(admin.AddDesigner))
			http.HandleFunc("/admin/files/edit", admin.CheckAuthCookie(admin.EditDesigner))

			http.Handle("/admin/static/", http.StripPrefix("/admin/static", http.FileServer(http.Dir("./admin/static"))))
		}
		if configuration.ServeWebsite {
			http.HandleFunc("/api/font", http2.GetFontMeta)
			http.HandleFunc("/", http2.GetWebApp)
		}

		log.Println("Serving at localhost:8090...")
		err = http.ListenAndServe(":8090", nil)
		if err != nil {
			panic(err)
		}
	}
}
