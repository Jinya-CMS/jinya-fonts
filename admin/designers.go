package admin

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
	"strings"
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

func DeleteDesigner(w http.ResponseWriter, r *http.Request) {
	fontName := r.URL.Query().Get("font")
	data, err := ioutil.ReadFile(config.LoadedConfiguration.FontFileFolder + "/" + fontName + ".yaml")
	if err != nil {
		RenderAdmin(w, "designers/delete", struct {
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
			Message      string
			FontName     string
			DesignerName string
		}{"Font does not exist", fontName, ""})
		return
	}

	if font.GoogleFont {
		RenderAdmin(w, "designers/index", struct {
			Message      string
			FontName     string
			DesignerName string
		}{"You cannot edit the designers of a synced font", fontName, ""})
		return
	}

	designerName := r.URL.Query().Get("name")
	var designer *meta.FontDesigner
	for _, item := range font.Designers {
		if item.Name == designerName {
			designer = &item
			break
		}
	}

	if designer == nil {
		RenderAdmin(w, "designers/index", struct {
			Message      string
			FontName     string
			DesignerName string
		}{"Designer doesn't exist in font", fontName, ""})
		return
	}

	if r.Method == http.MethodGet {
		RenderAdmin(w, "designers/delete", struct {
			Message      string
			FontName     string
			DesignerName string
		}{"", fontName, designerName})
	} else if r.Method == http.MethodPost {
		var designers []meta.FontDesigner
		for _, item := range font.Designers {
			if item.Name != designerName {
				designers = append(designers, item)
			}
		}

		font.Designers = designers

		data, err := yaml.Marshal(font)
		if err != nil {
			RenderAdmin(w, "designers/delete", struct {
				Message      string
				FontName     string
				DesignerName string
			}{"Failed to remove designer from font", fontName, designerName})
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+fontName+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "designers/delete", struct {
				Message      string
				FontName     string
				DesignerName string
			}{"Failed to remove designer from font", fontName, designerName})
			return
		}

		http.Redirect(w, r, "/admin/designers?font="+fontName, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AddDesigner(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == http.MethodGet {
		RenderAdmin(w, "designers/add", struct {
			Message  string
			FontName string
			Name     string
			Bio      string
		}{"", fontName, "", ""})
	} else if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			RenderAdmin(w, "designers/add", struct {
				Message  string
				FontName string
				Name     string
				Bio      string
			}{"Error creating font", fontName, "", ""})
			return
		}

		designerName := r.FormValue("name")
		designerBio := r.FormValue("bio")

		designerExists := false
		for _, item := range font.Designers {
			if strings.ToLower(item.Name) == strings.ToLower(designerName) {
				designerExists = true
				break
			}
		}

		if designerExists {
			RenderAdmin(w, "designers/add", struct {
				Message  string
				FontName string
				Name     string
				Bio      string
			}{"A designer with the chosen name is already present", fontName, designerName, designerBio})
			return
		}

		font.Designers = append(font.Designers, meta.FontDesigner{
			Name: designerName,
			Bio:  designerBio,
		})
		data, err := yaml.Marshal(font)
		if err != nil {
			RenderAdmin(w, "designers/add", struct {
				Message  string
				FontName string
				Name     string
				Bio      string
			}{"Failed to add designer to font", fontName, designerName, designerBio})
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+fontName+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "designers/add", struct {
				Message  string
				FontName string
				Name     string
				Bio      string
			}{"Failed to add designer to font", fontName, designerName, designerBio})
			return
		}

		http.Redirect(w, r, "/admin/designers?font="+fontName, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func EditDesigner(w http.ResponseWriter, r *http.Request) {
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

	designerNameFromUri := r.URL.Query().Get("name")
	var designer *meta.FontDesigner
	for _, item := range font.Designers {
		if item.Name == designerNameFromUri {
			designer = &item
			break
		}
	}

	if designer == nil {
		RenderAdmin(w, "designers/index", struct {
			Message   string
			FontName  string
			Designers []meta.FontDesigner
		}{"Designer doesn't exist in font", fontName, []meta.FontDesigner{}})
		return
	}

	if r.Method == http.MethodGet {
		RenderAdmin(w, "designers/edit", struct {
			Message  string
			FontName string
			Name     string
			Bio      string
		}{"", fontName, designer.Name, designer.Bio})
	} else if r.Method == http.MethodPost {
		err = r.ParseForm()
		if err != nil {
			RenderAdmin(w, "designers/edit", struct {
				Message  string
				FontName string
				Name     string
				Bio      string
			}{"Error creating font", fontName, "", ""})
			return
		}

		designerName := r.FormValue("name")
		designerBio := r.FormValue("bio")

		designerExists := false
		for _, item := range font.Designers {
			if strings.ToLower(item.Name) == strings.ToLower(designerName) && strings.ToLower(designerNameFromUri) != strings.ToLower(item.Name) {
				designerExists = true
				break
			}
		}

		if designerExists {
			RenderAdmin(w, "designers/edit", struct {
				Message  string
				FontName string
				Name     string
				Bio      string
			}{"A designer with the chosen name is already present", fontName, designerName, designerBio})
			return
		}

		for idx, item := range font.Designers {
			if strings.ToLower(item.Name) == strings.ToLower(designerNameFromUri) {
				font.Designers[idx].Name = designerName
				font.Designers[idx].Bio = designerBio
			}
		}

		data, err := yaml.Marshal(font)
		if err != nil {
			RenderAdmin(w, "designers/edit", struct {
				Message  string
				FontName string
				Name     string
				Bio      string
			}{"Failed to add designer to font", fontName, designerName, designerBio})
		}

		err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+fontName+".yaml", data, 0775)
		if err != nil {
			RenderAdmin(w, "designers/edit", struct {
				Message  string
				FontName string
				Name     string
				Bio      string
			}{"Failed to add designer to font", fontName, designerName, designerBio})
			return
		}

		http.Redirect(w, r, "/admin/designers?font="+fontName, http.StatusFound)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
