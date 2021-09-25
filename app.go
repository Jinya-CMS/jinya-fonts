package main

import (
	"flag"
	"jinya-fonts/config"
	"jinya-fonts/fontsync"
	http2 "jinya-fonts/http"
	"log"
	"net/http"
	"os"
)

func ContainsString(slice []string, search string) bool {
	for _, item := range slice {
		if item == search {
			return true
		}
	}

	return false
}

func main() {
	configFileFlag := flag.String("config-file", "./config.yaml", "The config file, check the sample for the structure")
	flag.Parse()

	configuration, err := config.LoadConfiguration(*configFileFlag)
	if err != nil {
		panic(err)
	}

	if ContainsString(os.Args, "sync") {
		err = fontsync.Sync(configuration)
		if err != nil {
			panic(err)
		}
	}

	if ContainsString(os.Args, "serve") {
		http.HandleFunc("/fonts/", http2.GetFont)
		http.HandleFunc("/css2", http2.GetCss2)
		log.Println("Serving at localhost:8090...")
		err = http.ListenAndServe(":8090", nil)
		if err != nil {
			panic(err)
		}
	}
}
