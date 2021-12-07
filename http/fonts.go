package http

import (
	"io"
	"jinya-fonts/config"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetFont(w http.ResponseWriter, r *http.Request) {
	configuration := config.LoadedConfiguration
	filename := strings.TrimPrefix(r.URL.Path, "/fonts/")
	path := configuration.FontFileFolder + "/" + filename
	file, err := os.Open(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	filestat, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "font/woff2")
	w.Header().Set("Content-Length", strconv.Itoa(int(filestat.Size())))
	_, err = io.Copy(w, file)
	if err != nil {
		http.NotFound(w, r)
		return
	}
}
