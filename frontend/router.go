package frontend

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/text/language"
	"html/template"
	"jinya-fonts/database"
	"net/http"
	"strings"
	"time"
)

//go:embed tmpl
var tmplFs embed.FS

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

func funcMap(w http.ResponseWriter, r *http.Request) template.FuncMap {
	return map[string]any{
		"year": getYear,
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"translate": translate(w, r),
	}
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("layout").Funcs(funcMap(w, r)).ParseFS(tmplFs, "tmpl/layout.gohtml", "tmpl/index.gohtml")
	if err == nil {
		t.Execute(w, nil)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func detailPage(w http.ResponseWriter, r *http.Request) {
	font, err := database.GetFont(r.URL.Query().Get("font"))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	designers := make([]string, len(font.Designers))
	for i, d := range font.Designers {
		designers[i] = d.Name
	}

	t, err := template.New("layout").Funcs(funcMap(w, r)).ParseFS(tmplFs, "tmpl/layout.gohtml", "tmpl/detail.gohtml")
	if err == nil {
		t.Execute(w, struct {
			Font      *database.Webfont
			Designers string
		}{
			Font:      font,
			Designers: strings.Join(designers, ", "),
		})
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func SetupRouter(router *mux.Router) {
	router.Methods(http.MethodGet).Path("/font").HandlerFunc(detailPage)
	router.Methods(http.MethodGet).PathPrefix("/").HandlerFunc(indexPage)
}
