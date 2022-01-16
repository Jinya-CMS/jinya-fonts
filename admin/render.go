package admin

import (
	"html/template"
	"net/http"
)

func RenderAdmin(w http.ResponseWriter, path string, data interface{}) error {
	tmpl, err := template.New("layout").ParseFiles("./admin/layout.gohtml", "./admin/"+path+".gohtml")
	if err != nil {
		return err
	}

	return tmpl.ExecuteTemplate(w, "layout", data)
}
