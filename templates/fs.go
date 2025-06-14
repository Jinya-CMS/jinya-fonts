package templates

import (
	"embed"
	fs2 "io/fs"
	"jinya-fonts/config"
	"os"
)

//go:embed admin
var adminTmpl embed.FS

//go:embed frontend
var frontendTmpl embed.FS

func getDevTemplatesFs() fs2.FS {
	return os.DirFS("templates")
}

func GetAdminTemplatesFs() fs2.FS {
	if config.IsDev() {
		return getDevTemplatesFs()
	}

	return adminTmpl
}

func GetFrontendTemplatesFs() fs2.FS {
	if config.IsDev() {
		return getDevTemplatesFs()
	}

	return frontendTmpl
}
