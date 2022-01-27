package http

import (
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
	"os"
)

func checkAuthCookie(r *http.Request) bool {
	authCookie, err := r.Cookie("auth")

	return err != nil || authCookie.Value != config.LoadedConfiguration.AdminPassword
}

func AdminFontConfig(w http.ResponseWriter, r *http.Request) {
	if checkAuthCookie(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := r.PostFormValue("name")
	if name == "" {
		w.Write([]byte("name"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	category := r.PostFormValue("category")
	if category == "" {
		w.Write([]byte("category"))
		w.WriteHeader(http.StatusConflict)
		return
	}

	license := r.PostFormValue("license")
	description := r.PostFormValue("description")

	font := meta.FontFile{
		Name:        name,
		Fonts:       []meta.FontFileMeta{},
		Description: description,
		Designers:   []meta.FontDesigner{},
		License:     license,
		Category:    category,
	}

	data, err := yaml.Marshal(font)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+name+".yaml", data, 0775)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = os.MkdirAll(config.LoadedConfiguration.FontFileFolder+"/"+name, 0775)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func AdminUploadFont(w http.ResponseWriter, r *http.Request) {
	if checkAuthCookie(r) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 * 1024 * 1024 * 1024)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	name := r.URL.Query().Get("name")
	file, err := os.Open(config.LoadedConfiguration.FontFileFolder + "/" + name + ".yaml")
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var metaData meta.FontFile
	err = yaml.Unmarshal(data, &metaData)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	uploadedFile, _, err := r.FormFile("fontfile")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	subset := r.FormValue("subset")
	variant := r.FormValue("variant")
	unicodeRange := r.FormValue("unicodeRange")
	weight := r.FormValue("weight")
	style := r.FormValue("style")

	path := config.LoadedConfiguration.FontFileFolder + "/" + name + "/" + name + "." + subset + "." + variant + ".woff2"
	fontData := meta.FontFileMeta{
		Path:         path,
		Subset:       subset,
		UnicodeRange: unicodeRange,
		Weight:       weight,
		Style:        style,
		FontName:     name,
	}

	targetFile, err := os.Open(path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer targetFile.Close()
	_, err = io.Copy(targetFile, uploadedFile)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	metaData.Fonts = append(metaData.Fonts, fontData)
	marshalledData, err := yaml.Marshal(metaData)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	err = ioutil.WriteFile(config.LoadedConfiguration.FontFileFolder+"/"+name+".yaml", marshalledData, 0775)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
