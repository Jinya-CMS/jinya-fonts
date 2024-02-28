package http

import (
	"archive/zip"
	"bytes"
	"github.com/gosimple/slug"
	"io"
	"jinya-fonts/config"
	"jinya-fonts/database"
	"net/http"
	"os"
)

func DownloadFont(w http.ResponseWriter, r *http.Request) {
	font := r.URL.Query().Get("font")

	if font == "" {
		http.NotFound(w, r)
		return
	}

	_, err := os.Stat(config.LoadedConfiguration.FontFileFolder + "/" + font)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	webfont, err := database.GetFont(font)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	css := ""
	files := make(map[string]*os.File, len(webfont.Fonts))

	for _, item := range webfont.Fonts {
		convertedCss, err := convertTemplateDataToCss(templateData{
			Name:        webfont.Name,
			Style:       item.Style,
			Url:         "./" + slug.Make(webfont.Name) + "/" + item.Path,
			Weight:      item.Weight,
			FontDisplay: "block",
		})
		if err != nil {
			continue
		}

		file, err := os.Open(config.LoadedConfiguration.FontFileFolder + "/" + webfont.Name + "/" + item.Path)
		if err != nil {
			continue
		}

		css += convertedCss + "\n"
		files[item.Path] = file
	}

	buffer := bytes.NewBuffer([]byte{})
	zipWriter := zip.NewWriter(buffer)

	for path, file := range files {
		_, err := file.Stat()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		zipHeaderWriter, err := zipWriter.Create(slug.Make(font) + "/" + path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(zipHeaderWriter, file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	zipFileWriter, err := zipWriter.Create(slug.Make(font) + ".css")
	_, err = zipFileWriter.Write([]byte(css))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = zipWriter.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+slug.Make(font)+".zip\"")
	_, _ = w.Write(buffer.Bytes())
}
