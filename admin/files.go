package admin

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
	"sort"
	"strings"
)

func FilesIndex(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font")
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"Font does not exist", fontName, []meta.FontFileMeta{}})
		return
	}

	var font meta.FontFile
	err = yaml.Unmarshal(data, &font)
	if err != nil {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"Font does not exist", fontName, []meta.FontFileMeta{}})
		return
	}

	if font.GoogleFont {
		RenderAdmin(w, "files/index", struct {
			Message  string
			FontName string
			Files    []meta.FontFileMeta
		}{"You cannot edit the files of a synced font", fontName, []meta.FontFileMeta{}})
		return
	}

	sort.Slice(font.Fonts, func(i, j int) bool {
		first := font.Fonts[i]
		second := font.Fonts[j]

		return strings.ToLower(first.Path) < strings.ToLower(second.Path)
	})
	err = RenderAdmin(w, "files/index", struct {
		Message  string
		FontName string
		Files    []meta.FontFileMeta
	}{"", font.Name, font.Fonts})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
