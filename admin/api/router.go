package api

import (
	"context"
	"crypto/rand"
	"github.com/gorilla/mux"
	"jinya-fonts/config"
	"net/http"

	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
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

	authN, err := authentication.New(ctx, zitadel.New(config.LoadedConfiguration.OpenIDDomain), string(encryptionKey), openid.DefaultAuthentication(config.LoadedConfiguration.OpenIDClientId, config.LoadedConfiguration.GetRedirectUrl(), string(encryptionKey)))
	if err != nil {
		panic(err)
	}

	mw := authentication.Middleware(authN)

	router.Methods("GET").Path("/api/admin/font/all").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(getAllFonts))))
	router.Methods("GET").Path("/api/admin/font/google").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(getGoogleFonts))))
	router.Methods("POST").Path("/api/admin/font/google").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(syncFonts))))
	router.Methods("GET").Path("/api/admin/font/custom").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(getCustomFonts))))
	router.Methods("POST").Path("/api/admin/font").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(createFont))))
	router.Methods("GET").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(getFontByName))))
	router.Methods("PUT").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(updateFont))))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(deleteFont))))

	router.Methods("GET").Path("/api/admin/font/{fontName}/file").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(getFontFiles))))
	router.Methods("POST").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(createFontFile))))
	router.Methods("PUT").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(updateFontFile))))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(deleteFontFile))))

	router.Methods("GET").Path("/api/admin/font/{fontName}/designer").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(getFontDesigners))))
	router.Methods("POST").Path("/api/admin/font/{fontName}/designer").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(createFontDesigner))))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}/designer/{designerName}").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(deleteFontDesigner))))

	router.Methods("GET").Path("/api/admin/settings").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(getSettings))))
	router.Methods("PUT").Path("/api/admin/settings").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(updateSettings))))

	router.Methods("GET").Path("/api/admin/status").Handler(mw.RequireAuthentication()(contentTypeJson()(http.HandlerFunc(getStatus))))

	router.Methods("GET").Path("/api/healthz").HandlerFunc(getHealth)
}
