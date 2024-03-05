package fonts

import (
	"archive/zip"
	"bytes"
	"github.com/gosimple/slug"
	"jinya-fonts/database"
	"net/http"
)

func downloadFont(w http.ResponseWriter, r *http.Request) {
	font := r.URL.Query().Get("font")

	if font == "" {
		http.NotFound(w, r)
		return
	}

	webfont, err := database.GetFont(font)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	zipWriter := zip.NewWriter(buffer)
	zipCssFileWriter, err := zipWriter.Create(slug.Make(font) + ".css")

	for _, item := range webfont.Fonts {
		convertedCss, err := convertTemplateDataToCss(templateData{
			File:        &item,
			Webfont:     webfont,
			FontDisplay: "block",
		})
		if err != nil {
			continue
		}

		fileName := database.GetFontFileName(webfont.Name, item.Weight, item.Style, item.Type, webfont.GoogleFont)
		file, _, err := getFontFile(fileName)
		if err != nil {
			return
		}
		if err != nil {
			continue
		}

		_, err = zipCssFileWriter.Write([]byte(convertedCss))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		zipFontFileWriter, err := zipWriter.Create(fileName)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = zipFontFileWriter.Write(file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Disposition", "attachment; filename=\""+slug.Make(font)+".zip\"")
	_, _ = w.Write(buffer.Bytes())
}
