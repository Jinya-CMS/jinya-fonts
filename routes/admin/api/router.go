package api

import (
	"context"
	"jinya-fonts/config"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

func contentTypeJson() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, req)
		})
	}
}

func SetupAdminApiRouter(router *mux.Router) {
	ctx := context.Background()

	zitadelConfig := oauth.WithIntrospection[*oauth.IntrospectionContext](oauth.ClientIDSecretIntrospectionAuthentication(config.LoadedConfiguration.OidcServerClientId, config.LoadedConfiguration.OidcServerClientSecret))
	authZ, err := authorization.New(ctx, zitadel.New(config.LoadedConfiguration.OidcDomain), zitadelConfig)

	if err != nil {
		panic(err)
	}

	mw := middleware.New(authZ)

	router.
		Methods(http.MethodGet).
		Path("/api/admin/font/all").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(GetAllFonts))))
	router.
		Methods(http.MethodGet).
		Path("/api/admin/font/google").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getGoogleFonts))))
	router.
		Methods(http.MethodPost).
		Path("/api/admin/font/google").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(syncFonts))))
	router.
		Methods(http.MethodGet).
		Path("/api/admin/font/custom").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getCustomFonts))))
	router.
		Methods(http.MethodPost).
		Path("/api/admin/font").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createFont))))
	router.
		Methods(http.MethodGet).
		Path("/api/admin/font/{fontName}").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(GetFontByName))))
	router.
		Methods(http.MethodPut).
		Path("/api/admin/font/{fontName}").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(updateFont))))
	router.
		Methods(http.MethodDelete).
		Path("/api/admin/font/{fontName}").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteFont))))

	router.
		Methods(http.MethodGet).
		Path("/api/admin/font/{fontName}/file").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(GetFontFiles))))
	router.
		Methods(http.MethodPost).
		Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createFontFile))))
	router.
		Methods(http.MethodPut).
		Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(updateFontFile))))
	router.
		Methods(http.MethodDelete).
		Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteFontFile))))

	router.
		Methods(http.MethodGet).
		Path("/api/admin/font/{fontName}/designer").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(GetFontDesigners))))
	router.
		Methods(http.MethodPost).
		Path("/api/admin/font/{fontName}/designer").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createFontDesigner))))
	router.
		Methods(http.MethodDelete).
		Path("/api/admin/font/{fontName}/designer/{designerName}").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteFontDesigner))))

	router.
		Methods(http.MethodGet).
		Path("/api/admin/settings").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getSettings))))
	router.
		Methods(http.MethodPut).
		Path("/api/admin/settings").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(updateSettings))))

	router.
		Methods(http.MethodGet).
		Path("/api/admin/status").
		Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getStatus))))

	router.
		Methods(http.MethodGet).
		Path("/api/healthz").
		HandlerFunc(getHealth)
}
