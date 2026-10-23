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

	adminRouter := router.
		PathPrefix("/api/admin").
		Subrouter()

	adminRouter.
		Methods(http.MethodGet).
		Path("/font/all").
		Handler(http.HandlerFunc(GetAllFonts))
	adminRouter.
		Methods(http.MethodGet).
		Path("/font/google").
		Handler(http.HandlerFunc(getGoogleFonts))
	adminRouter.
		Methods(http.MethodPost).
		Path("/font/google").
		Handler(http.HandlerFunc(syncFonts))
	adminRouter.
		Methods(http.MethodGet).
		Path("/font/custom").
		Handler(http.HandlerFunc(getCustomFonts))
	adminRouter.
		Methods(http.MethodPost).
		Path("/font").
		Handler(http.HandlerFunc(createFont))
	adminRouter.
		Methods(http.MethodGet).
		Path("/font/{fontName}").
		Handler(http.HandlerFunc(GetFontByName))
	adminRouter.
		Methods(http.MethodPut).
		Path("/font/{fontName}").
		Handler(http.HandlerFunc(updateFont))
	adminRouter.
		Methods(http.MethodDelete).
		Path("/font/{fontName}").
		Handler(http.HandlerFunc(deleteFont))

	adminRouter.
		Methods(http.MethodGet).
		Path("/font/{fontName}/file").
		Handler(http.HandlerFunc(GetFontFiles))
	adminRouter.
		Methods(http.MethodPost).
		Path("/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").
		Handler(http.HandlerFunc(createFontFile))
	adminRouter.
		Methods(http.MethodPut).
		Path("/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").
		Handler(http.HandlerFunc(updateFontFile))
	adminRouter.
		Methods(http.MethodDelete).
		Path("/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").
		Handler(http.HandlerFunc(deleteFontFile))

	adminRouter.
		Methods(http.MethodGet).
		Path("/font/{fontName}/designer").
		Handler(http.HandlerFunc(GetFontDesigners))
	adminRouter.
		Methods(http.MethodPost).
		Path("/font/{fontName}/designer").
		Handler(http.HandlerFunc(createFontDesigner))
	adminRouter.
		Methods(http.MethodDelete).
		Path("/font/{fontName}/designer/{designerName}").
		Handler(http.HandlerFunc(deleteFontDesigner))

	adminRouter.
		Methods(http.MethodGet).
		Path("/settings").
		Handler(http.HandlerFunc(getSettings))
	adminRouter.
		Methods(http.MethodPut).
		Path("/settings").
		Handler(http.HandlerFunc(updateSettings))

	adminRouter.
		Methods(http.MethodGet).
		Path("/status").
		Handler(http.HandlerFunc(getStatus))

	adminRouter.Use(mw.RequireAuthorization(), contentTypeJson())

	router.
		Methods(http.MethodGet).
		Path("/api/healthz").
		HandlerFunc(getHealth)
}
