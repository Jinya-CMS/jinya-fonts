package web

import (
	"html/template"
	"jinya-fonts/templates"
	"net/http"
	"os"
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("content").ParseFS(templates.GetAdminTemplatesFs(), "admin/index.gohtml")
	if err == nil {
		t.Execute(w, struct {
			OidcFrontendClientId string
			OidcDomain           string
			GetRedirectUrl       string
			ServerUrl            string
		}{
			OidcFrontendClientId: os.Getenv("OIDC_FRONTEND_CLIENT_ID"),
			OidcDomain:           os.Getenv("OIDC_DOMAIN"),
			ServerUrl:            os.Getenv("SERVER_URL"),
		})
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
