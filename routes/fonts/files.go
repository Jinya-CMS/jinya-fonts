package fonts

import (
	"bytes"
	"io"
	"jinya-fonts/storage"
	"net/http"
	"strings"
)

func getFontFile(path string) (stream io.Reader, fontType string, err error) {
	data, err := storage.GetFontFile(path)
	if err != nil {
		return nil, "", err
	}

	stream = bytes.NewBuffer(data)
	if strings.HasSuffix(path, "woff2") {
		fontType = "font/woff2"
	} else if strings.HasSuffix(path, "ttf") {
		fontType = "font/ttf"
	} else {
		fontType = "application/octet-stream"
	}

	return
}

type fileHandler struct{}

func (h fileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	splitPath := strings.Split(r.URL.Path, "/")
	path := splitPath[len(splitPath)-1]

	fontFile, fileType, err := getFontFile(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", fileType)
	w.WriteHeader(http.StatusOK)
	_, _ = io.Copy(w, fontFile)
}
