package admin

import (
	"embed"
	"html/template"
	"net/http"
)

var (
	//go:embed templates
	templates embed.FS
)

func RenderAdmin(w http.ResponseWriter, path string, data interface{}) error {
	tmpl, err := template.New("layout").ParseFS(templates, "templates/layout.gohtml", "templates/"+path+".gohtml")
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}
