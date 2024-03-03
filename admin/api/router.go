package api

import (
	"context"
	"github.com/gorilla/mux"
	"jinya-fonts/config"
	"net/http"

	"github.com/zitadel/zitadel-go/v3/pkg/authentication"
	openid "github.com/zitadel/zitadel-go/v3/pkg/authentication/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

func SetupAdminApiRouter(router *mux.Router) {
	ctx := context.Background()
	authN, err := authentication.New(ctx, zitadel.New(config.LoadedConfiguration.OpenIDDomain), config.LoadedConfiguration.EncryptionKey,
		openid.DefaultAuthentication(config.LoadedConfiguration.OpenIDClientId, config.LoadedConfiguration.GetRedirectUrl(), config.LoadedConfiguration.EncryptionKey),
	)
	if err != nil {
		panic(err)
	}

	mw := authentication.Middleware(authN)

	router.Methods("GET").Path("/api/admin/font/all").Handler(mw.RequireAuthentication()(http.HandlerFunc(getAllFonts)))
	router.Methods("GET").Path("/api/admin/font/google").Handler(mw.RequireAuthentication()(http.HandlerFunc(getGoogleFonts)))
	router.Methods("GET").Path("/api/admin/font/custom").Handler(mw.RequireAuthentication()(http.HandlerFunc(getCustomFonts)))
	router.Methods("POST").Path("/api/admin/font").Handler(mw.RequireAuthentication()(http.HandlerFunc(createFont)))
	router.Methods("GET").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthentication()(http.HandlerFunc(getFontByName)))
	router.Methods("PUT").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthentication()(http.HandlerFunc(updateFont)))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}").Handler(mw.RequireAuthentication()(http.HandlerFunc(deleteFont)))

	router.Methods("GET").Path("/api/admin/font/{fontName}/file").Handler(mw.RequireAuthentication()(http.HandlerFunc(getFontFiles)))
	router.Methods("POST").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthentication()(http.HandlerFunc(createFontFile)))
	router.Methods("PUT").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthentication()(http.HandlerFunc(updateFontFile)))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}/file/{fontWeight}.{fontStyle}.{fontType}").Handler(mw.RequireAuthentication()(http.HandlerFunc(deleteFontFile)))

	router.Methods("GET").Path("/api/admin/font/{fontName}/designer").Handler(mw.RequireAuthentication()(http.HandlerFunc(getFontDesigners)))
	router.Methods("POST").Path("/api/admin/font/{fontName}/designer").Handler(mw.RequireAuthentication()(http.HandlerFunc(createFontDesigner)))
	router.Methods("DELETE").Path("/api/admin/font/{fontName}/designer/{designerName}").Handler(mw.RequireAuthentication()(http.HandlerFunc(deleteFontDesigner)))

	router.Methods("GET").Path("/api/admin/settings").Handler(mw.RequireAuthentication()(http.HandlerFunc(getSettings)))
	router.Methods("PUT").Path("/api/admin/settings").Handler(mw.RequireAuthentication()(http.HandlerFunc(updateSettings)))
}
