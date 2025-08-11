package frontend

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"jinya-fonts/database"
	"jinya-fonts/storage"
	"jinya-fonts/templates"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/text/language"
)

//go:embed langs
var langsFs embed.FS

var supportedLanguages = []language.Tag{language.English, language.German}

func getYear() string {
	return time.Now().Format("2006")
}

func translate(w http.ResponseWriter, r *http.Request) func(key string, replacements ...any) string {
	errorFunc := func(key string, replacements ...any) string {
		return key
	}

	matcher := language.NewMatcher(supportedLanguages)
	acceptLanguage := r.Header.Get("Accept-Language")
	tag, _ := language.MatchStrings(matcher, acceptLanguage)
	base, _, _ := tag.Raw()
	messagesToUseBytes, err := langsFs.ReadFile("langs/messages." + base.String() + ".json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return errorFunc
	}

	var messages map[string]string
	err = json.Unmarshal(messagesToUseBytes, &messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return errorFunc
	}

	return func(key string, replacements ...any) string {
		translation, ok := messages[key]
		if !ok {
			return key
		}

		return fmt.Sprintf(translation, replacements...)
	}
}

func render(w http.ResponseWriter, r *http.Request, name string, data any) {
	t, err := template.New("layout").Funcs(map[string]any{
		"year": getYear,
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"translate": translate(w, r),
	}).ParseFS(templates.GetFrontendTemplatesFs(), "frontend/layout.gohtml", fmt.Sprintf("frontend/%s.gohtml", name))
	if err == nil {
		t.Execute(w, data)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	render(w, r, "index", nil)
}

func detailPage(w http.ResponseWriter, r *http.Request) {
	font, err := database.Get[database.Webfont](r.URL.Query().Get("font"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	designers, err := database.GetDesigners(font.Name)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	files, err := storage.GetFontFiles(font.Name)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	font.Fonts = files

	des := make([]string, len(designers))
	for i, d := range designers {
		des[i] = d.Name
	}

	render(w, r, "detail", struct {
		Font      *database.Webfont
		Designers string
	}{
		Font:      font,
		Designers: strings.Join(des, ", "),
	})
}

func SetupRouter(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/font").HandlerFunc(detailPage)
	router.Methods(http.MethodGet).PathPrefix("/").HandlerFunc(indexPage)
}
