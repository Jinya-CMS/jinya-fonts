package http

import (
	"io"
	"jinya-fonts/config"
	"net/http"
	"os"
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

	_, err = io.Copy(w, file)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}