package http

import (
	"gopkg.in/yaml.v3"
	"html/template"
	"io"
	"io/ioutil"
	"jinya-fonts/config"
	"jinya-fonts/meta"
	"net/http"
	"os"
	"strings"
	"time"
)

func checkAuthCookie(r *http.Request) bool {
	authCookie, err := r.Cookie("auth")

	return err != nil || authCookie.Value != config.LoadedConfiguration.AdminPassword
}

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		path := strings.TrimPrefix(r.URL.Path, "/login")
		if path == "/" || path == "" {
			path = "login.html"
		}

		data, err := ioutil.ReadFile("./admin/" + path)
		if err != nil {
			path = "login.html"
			data, err = ioutil.ReadFile("./admin/" + path)

			if err != nil {
				http.NotFound(w, r)
				return
			}
		}

		if strings.HasSuffix(path, "css") {
			w.Header().Set("Content-Type", "text/css")
		} else if strings.HasSuffix(path, "html") {
			w.Header().Set("Content-Type", "text/html")
		} else if strings.HasSuffix(path, "js") {
			w.Header().Set("Content-Type", "application/javascript")
		}

		w.Write(data)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		password := r.FormValue("password")
		if password == config.LoadedConfiguration.AdminPassword {
			cookie := &http.Cookie{
				Name:     "auth",
				Value:    password,
				Path:     "/",
				Expires:  time.Now().Add(time.Hour * 24),
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
			}

			remember := r.FormValue("remember") == "on"
			if remember {
				cookie.Expires = time.Date(9999, 12, 31, 23, 59, 59, 0, time.UTC)
			}
			http.SetCookie(w, cookie)
			http.Redirect(w, r, "/admin", http.StatusFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	if checkAuthCookie(r) {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	type fontData struct {
		Name         string
		NumberStyles int
		License      string
		Category     string
		Author       string
		GoogleFont   bool
	}

	tmpl, err := template.ParseFiles("./admin/index.gohtml")
	if err != nil {
		http.NotFound(w, r)
		return
	}

	fonts, err := loadFonts()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var data []fontData
	for _, font := range fonts {
		var designers []string
		for _, designer := range font.Designers {
			designers = append(designers, designer.Name)
		}

		data = append(data, fontData{
			Name:         font.Name,
			NumberStyles: len(font.Fonts),
			License:      font.License,
			Category:     font.Category,
			Author:       strings.Join(designers, ","),
			GoogleFont:   font.GoogleFont,
		})
	}

	tmpl.Execute(w, data)
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
		Variant:      variant,
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
