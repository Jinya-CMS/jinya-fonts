package admin

import (
	"bytes"
	"jinya-fonts/config"
	"jinya-fonts/fontsync"
	"log"
	"net/http"
)

func TriggerSync(w http.ResponseWriter, r *http.Request) {
	type syncData struct {
		Log string
	}

	if r.Method == http.MethodPost {
		backupWriter := log.Writer()
		buffer := bytes.Buffer{}
		log.SetOutput(&buffer)

		_ = fontsync.Sync(config.LoadedConfiguration)

		log.SetOutput(backupWriter)

		data := syncData{Log: buffer.String()}
		err := RenderAdmin(w, "triggerSync", data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		err := RenderAdmin(w, "triggerSync", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
