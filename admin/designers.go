package admin

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
)

func DesignersIndex(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font")
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
	if err != nil {
		RenderAdmin(w, "designers/index", struct {
			Message   string
			FontName  string
			Designers []meta.FontDesigner
		}{"Font does not exist", fontName, []meta.FontDesigner{}})
		return
	}

	var font meta.FontFile
	err = yaml.Unmarshal(data, &font)
	if err != nil {
		RenderAdmin(w, "designers/index", struct {
			Message   string
			FontName  string
			Designers []meta.FontDesigner
		}{"Font does not exist", fontName, []meta.FontDesigner{}})
		return
	}

	if font.GoogleFont {
		RenderAdmin(w, "designers/index", struct {
			Message   string
			FontName  string
			Designers []meta.FontDesigner
		}{"You cannot edit the designers of a synced font", fontName, []meta.FontDesigner{}})
		return
	}

	err = RenderAdmin(w, "designers/index", struct {
		Message   string
		FontName  string
		Designers []meta.FontDesigner
	}{"", font.Name, font.Designers})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
