package fonts

import (
	"archive/zip"
	"bytes"
	"github.com/gosimple/slug"
	"io"
	"jinya-fonts/database"
	"jinya-fonts/storage"
	"net/http"
)

func downloadFont(w http.ResponseWriter, r *http.Request) {
	font := r.URL.Query().Get("font")

	if font == "" {
		http.NotFound(w, r)
		return
	}

	webfont, err := database.Get[database.Webfont](font)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	files, err := storage.GetFontFiles(webfont.Name)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	zipWriter := zip.NewWriter(buffer)
	css := ""

	for _, item := range files {
		convertedCss, err := convertTemplateDataToCss(templateData{
			File:        &item,
			Webfont:     webfont,
			FontDisplay: "block",
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		css = css + "\n" + convertedCss

		fileName := storage.GetFontFileName(webfont.Name, item.Weight, item.Style, item.Type, webfont.GoogleFont)
		file, _, err := getFontFile(fileName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		zipFontFileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(zipFontFileWriter, file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	zipCssFileWriter, err := zipWriter.Create(slug.Make(font) + ".css")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = zipCssFileWriter.Write([]byte(css))
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
