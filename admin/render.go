package admin

import (
	"html/template"
	"net/http"
)

func RenderAdmin(w http.ResponseWriter, path string, data interface{}) error {
	tmpl, err := template.New("layout").ParseFiles("./admin/templates/layout.gohtml", "./admin/templates/"+path+".gohtml")
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}
