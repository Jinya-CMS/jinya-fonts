package fonts

import (
	"fmt"
	"jinya-fonts/database"
	"net/http"
	"strings"
)

func getMetadataFromFile(filename string) (name, weight, style, fileType string) {
	split := strings.Split(strings.ToLower(filename), ".")
	if len(split) < 5 {
		return "", "", "", ""
	}

	name = split[0]
	weight = split[1]
	style = split[2]
	fileType = split[3]

	return
}

func getFontFile(path string) (data []byte, fontType string, err error) {
	data, err = database.GetCachedFontFile(path)
	if err != nil {
		file, weight, style, fileType := getMetadataFromFile(path)
		var font *database.Webfont
		font, err = database.GetFont(file)
		if err != nil {
			return
		}

		fileFromDatabase, exists := font.Fonts[database.GetFontFileName(file, weight, style, fileType, font.GoogleFont)]
		if !exists {
			err = fmt.Errorf("not found")
			return
		}

		_ = database.AddCachedFontFile(file, weight, style, fileType, fileFromDatabase.Data, font.GoogleFont)

		data = fileFromDatabase.Data
		fontType = "font/" + fileFromDatabase.Type
	} else {
		if strings.HasSuffix(path, "woff2") {
			fontType = "font/woff2"
		} else if strings.HasSuffix(path, "ttf") {
			fontType = "font/ttf"
		} else {
			fontType = "application/octet-stream"
		}
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
	_, _ = w.Write(fontFile)
}
