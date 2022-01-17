package http

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func GetWebApp(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" || path == "" {
		path = "index.html"
	}

	data, err := ioutil.ReadFile("./webapp/" + path)
	if err != nil {
		path = "index.html"
		data, err = ioutil.ReadFile("./webapp/" + path)

		if err != nil {
			http.NotFound(w, r)
			return
		}
	}

	if strings.HasSuffix(path, "css") {
		w.Header().Set("Content-Type", "text/css")
	} else if strings.HasSuffix(path, "html") {
		w.Header().Set("Content-Type", "text/html")
	} else if strings.HasSuffix(path, "js") {
		w.Header().Set("Content-Type", "application/javascript")
	}

	w.Write(data)
}
