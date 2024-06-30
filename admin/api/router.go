package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"jinya-fonts/config"
	"net/http"

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
	encryptionKey := make([]byte, 32)
	_, err := rand.Read(encryptionKey)
	if err != nil {
		panic(err)
	}

	keyFileData, err := base64.StdEncoding.DecodeString(config.LoadedConfiguration.OpenIDKeyFileData)
	if err != nil {
		panic(err)
	}

	keyFile, err := client.ConfigFromKeyFileData(keyFileData)
	zitadelConfig := oauth.WithIntrospection[*oauth.IntrospectionContext](oauth.JWTProfileIntrospectionAuthentication(keyFile))
	authZ, err := authorization.New(ctx, zitadel.New(config.LoadedConfiguration.OpenIDDomain), zitadelConfig)

	if err != nil {
		panic(err)
	}

	mw := middleware.New(authZ)

	router.Methods("GET").Path("/api/admin/font/all").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getAllFonts))))
	router.Methods("GET").Path("/api/admin/font/google").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getGoogleFonts))))
	router.Methods("POST").Path("/api/admin/font/google").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(syncFonts))))
	router.Methods("GET").Path("/api/admin/font/custom").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getCustomFonts))))
	router.Methods("POST").Path("/api/admin/font").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createFont))))
	router.Methods("GET").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getFontByName))))
	router.Methods("PUT").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(updateFont))))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteFont))))

	router.Methods("GET").Path("/api/admin/font/{fontName}/file").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getFontFiles))))
	router.Methods("POST").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createFontFile))))
	router.Methods("PUT").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(updateFontFile))))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteFontFile))))

	router.Methods("GET").Path("/api/admin/font/{fontName}/designer").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getFontDesigners))))
	router.Methods("POST").Path("/api/admin/font/{fontName}/designer").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(createFontDesigner))))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}/designer/{designerName}").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(deleteFontDesigner))))

	router.Methods("GET").Path("/api/admin/settings").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getSettings))))
	router.Methods("PUT").Path("/api/admin/settings").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(updateSettings))))

	router.Methods("GET").Path("/api/admin/status").Handler(mw.RequireAuthorization()(contentTypeJson()(http.HandlerFunc(getStatus))))

	router.Methods("GET").Path("/api/healthz").HandlerFunc(getHealth)
}
